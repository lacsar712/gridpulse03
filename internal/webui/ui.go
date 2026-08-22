package webui

func IndexHTML() []byte {
	return []byte(`<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8"/>
  <meta name="viewport" content="width=device-width, initial-scale=1"/>
  <title>GridPulse Relay</title>
  <link rel="stylesheet" href="/static/app.css"/>
</head>
<body>
  <header>
    <h1>GridPulse Relay</h1>
    <p class="sub">电网继电保护跳闸 / 保护出口投递</p>
  </header>
  <main>
    <section class="panel">
      <h2>New record</h2>
      <form id="create-form">
        <label>Title <input name="title" required maxlength="120"/></label>
        <label>Tags <input name="tags" placeholder="ops,runbook"/></label>
        <label>Body <textarea name="body" rows="6" placeholder="see PROJECT.md"></textarea></label>
        <button type="submit">Create</button>
      </form>
    </section>
    <section class="panel">
      <h2>Search</h2>
      <input id="q" placeholder="filter title or body"/>
      <button id="reload" type="button">Reload</button>
      <div id="status"></div>
      <ul id="list"></ul>
    </section>
  </main>
  <script src="/static/app.js"></script>
</body>
</html>
`)
}

func AppCSS() []byte {
	return []byte(`body{font-family:Georgia,serif;margin:0;background:#f4f1ea;color:#222}
header{background:#243447;color:#fff;padding:24px 32px}
h1{margin:0;font-size:28px}
.sub{opacity:.8;margin:6px 0 0}
main{display:grid;grid-template-columns:1fr 1fr;gap:20px;padding:24px}
.panel{background:#fff;border:1px solid #ddd;padding:16px;border-radius:8px}
label{display:block;margin:8px 0}
input,textarea{width:100%;box-sizing:border-box;padding:8px}
button{background:#243447;color:#fff;border:0;padding:8px 14px;cursor:pointer}
#list{list-style:none;padding:0}
#list li{border-bottom:1px solid #eee;padding:8px 0}
#status{min-height:1.2em;color:#666;margin:8px 0}
@media (max-width:800px){main{grid-template-columns:1fr}}
`)
}

func AppJS() []byte {
	return []byte(`async function load(){
  const q = document.getElementById('q').value;
  const st = document.getElementById('status');
  st.textContent = 'loading...';
  const r = await fetch('/api/records?q=' + encodeURIComponent(q));
  const items = await r.json();
  const ul = document.getElementById('list');
  ul.innerHTML = '';
  (items||[]).forEach(it => {
    const li = document.createElement('li');
    li.textContent = '#' + it.id + ' ' + it.title + ' [' + (it.tags||[]).join(',') + ']';
    ul.appendChild(li);
  });
  st.textContent = (items||[]).length + ' records';
}
document.getElementById('reload').onclick = load;
document.getElementById('create-form').onsubmit = async (e) => {
  e.preventDefault();
  const fd = new FormData(e.target);
  const tags = (fd.get('tags')||'').toString().split(',').map(s=>s.trim()).filter(Boolean);
  const r = await fetch('/api/records', {
    method:'POST',
    headers:{'Content-Type':'application/json'},
    body: JSON.stringify({title:fd.get('title'), body:fd.get('body'), tags})
  });
  if(!r.ok){ document.getElementById('status').textContent = 'create failed'; return; }
  e.target.reset();
  load();
};
load();
`)
}
