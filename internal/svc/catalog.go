package svc

import (
	"context"
	"fmt"
	"strings"

	"example.com/gridpulse/internal/auth"
	"example.com/gridpulse/internal/clock"
	"example.com/gridpulse/internal/config"
	"example.com/gridpulse/internal/ids"
	"example.com/gridpulse/internal/model"
	"example.com/gridpulse/internal/notify"
	"example.com/gridpulse/internal/policy"
	"example.com/gridpulse/internal/store"
	"example.com/gridpulse/internal/util/errwrap"
	"example.com/gridpulse/internal/util/kvbag"
	"example.com/gridpulse/internal/util/pathsafe"
	"example.com/gridpulse/internal/util/precopy"
	"example.com/gridpulse/internal/validate"
)

type Catalog struct {
	Store  *store.Store
	Config config.Config
	Meta   *kvbag.Bag
	Clk    clock.Clock
	Events *notify.Sink
}

func NewCatalog(st *store.Store, cfg config.Config) *Catalog {
	return &Catalog{
		Store:  st,
		Config: cfg,
		Meta:   kvbag.New(),
		Clk:    clock.Real{},
		Events: &notify.Sink{},
	}
}

func (c *Catalog) BootstrapAdmin(ctx context.Context) error {
	if err := validate.Token(c.Config.AdminToken); err != nil {
		return err
	}
	users, err := c.Store.ListUsers(ctx)
	if err != nil {
		return err
	}
	now := clock.RFC3339(c.Clk.Now())
	if len(users) == 0 {
		_, err = c.Store.CreateUser(ctx, "admin", auth.HashToken(c.Config.AdminToken), "admin")
		return err
	}
	return c.Store.SetSetting(ctx, "bootstrapped", now)
}

func (c *Catalog) Create(ctx context.Context, title, body string, tags []string, owner int64) (model.Record, error) {
	title = strings.TrimSpace(title)
	if err := validate.Title(title); err != nil {
		c.Events.Error("reject title: " + err.Error())
		return model.Record{}, fmt.Errorf("create: %w", errwrap.ErrDenied)
	}
	if err := validate.Body(body, int(c.Config.MaxBodyBytes)); err != nil {
		return model.Record{}, err
	}
	if err := policy.Enforce(title, body, tags); err != nil {
		c.Events.Error("policy: " + err.Error())
		return model.Record{}, fmt.Errorf("create: %w", errwrap.WrapDenied("policy"))
	}
	if err := policy.AfterWrite(
		func() (string, error) { return c.Store.GetSetting(ctx, "rollback_min") },
		func(v string) error { return c.Store.SetSetting(ctx, "rollback_min", v) },
		body,
	); err != nil {
		c.Events.Error("after-write: " + err.Error())
		return model.Record{}, err
	}
	slug := model.Slugify(title)
	rec := model.Record{Slug: slug, Title: title, Body: body, Tags: tags, OwnerID: owner}
	out, err := c.Store.CreateRecord(ctx, rec)
	if err != nil {
		return model.Record{}, errwrap.Wrapf(err, "create")
	}
	c.Store.Remember(&out)
	_ = c.Store.AddAudit(ctx, "system", "create", ids.Key("rec", out.Slug, ids.ShortID(8)))
	c.Events.Info("created " + out.Slug)
	return out, nil
}

func (c *Catalog) Update(ctx context.Context, id int64, title, body string, tags []string, editor string) (model.Record, error) {
	if err := validate.Title(title); err != nil {
		return model.Record{}, err
	}
	if err := policy.Enforce(title, body, tags); err != nil {
		return model.Record{}, err
	}
	cur, err := c.Store.GetRecord(ctx, id)
	if err != nil {
		return model.Record{}, err
	}
	cur.Title = title
	cur.Body = body
	cur.Tags = tags
	out, err := c.Store.UpdateRecord(ctx, cur, editor)
	if err != nil {
		return out, err
	}
	c.Store.Remember(&out)
	return out, nil
}

func (c *Catalog) Get(ctx context.Context, id int64) (model.Record, error) {
	return c.Store.GetRecord(ctx, id)
}

func (c *Catalog) GetBySlug(ctx context.Context, slug string) (model.Record, error) {
	if rec, ok := c.Store.Cached(slug); ok {
		return *rec, nil
	}
	rec, err := c.Store.GetBySlug(ctx, slug)
	if err != nil {
		return model.Record{}, err
	}
	c.Store.Remember(&rec)
	return rec, nil
}

func (c *Catalog) List(ctx context.Context, f model.ListFilter) ([]model.Record, error) {
	f = f.Normalized(validate.Limit(f.Limit, c.Config.PageSize))
	return c.Store.ListRecords(ctx, f)
}

func (c *Catalog) PreviewTags(tags []string, n int) []string {
	return precopy.HeadStrings(tags, n)
}

func (c *Catalog) AttachPath(name string) (string, error) {
	safe := ids.ShortID(8) + "-" + name
	return pathsafe.JoinUnder(c.Config.DataDir, safe)
}

func (c *Catalog) AddFile(ctx context.Context, recordID int64, name string) (model.Attachment, error) {
	path, err := c.AttachPath(name)
	if err != nil {
		return model.Attachment{}, err
	}
	a := model.Attachment{RecordID: recordID, Name: name, SHA: ids.NewSlug("blob"), Path: path}
	return c.Store.AddAttachment(ctx, a)
}

func (c *Catalog) Detail(ctx context.Context, id int64) (model.Record, []model.Revision, []model.Attachment, error) {
	rec, err := c.Get(ctx, id)
	if err != nil {
		return model.Record{}, nil, nil, err
	}
	revs, err := c.Store.Revisions(ctx, id)
	if err != nil {
		return rec, nil, nil, err
	}
	atts, err := c.Store.Attachments(ctx, id)
	return rec, revs, atts, err
}
