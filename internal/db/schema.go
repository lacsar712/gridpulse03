package db

const schema = "" +
	"CREATE TABLE IF NOT EXISTS users (" +
	" id INTEGER PRIMARY KEY AUTOINCREMENT," +
	" name TEXT NOT NULL UNIQUE," +
	" token_hash TEXT NOT NULL," +
	" role TEXT NOT NULL," +
	" created_at TEXT NOT NULL" +
	");" +
	"CREATE TABLE IF NOT EXISTS records (" +
	" id INTEGER PRIMARY KEY AUTOINCREMENT," +
	" slug TEXT NOT NULL UNIQUE," +
	" title TEXT NOT NULL," +
	" body TEXT NOT NULL," +
	" tags TEXT NOT NULL," +
	" owner_id INTEGER NOT NULL," +
	" bytes INTEGER NOT NULL," +
	" created_at TEXT NOT NULL," +
	" updated_at TEXT NOT NULL," +
	" FOREIGN KEY(owner_id) REFERENCES users(id)" +
	");" +
	"CREATE TABLE IF NOT EXISTS revisions (" +
	" id INTEGER PRIMARY KEY AUTOINCREMENT," +
	" record_id INTEGER NOT NULL," +
	" body TEXT NOT NULL," +
	" editor TEXT NOT NULL," +
	" created_at TEXT NOT NULL," +
	" FOREIGN KEY(record_id) REFERENCES records(id)" +
	");" +
	"CREATE TABLE IF NOT EXISTS attachments (" +
	" id INTEGER PRIMARY KEY AUTOINCREMENT," +
	" record_id INTEGER NOT NULL," +
	" name TEXT NOT NULL," +
	" sha TEXT NOT NULL," +
	" size INTEGER NOT NULL," +
	" path TEXT NOT NULL," +
	" created_at TEXT NOT NULL," +
	" FOREIGN KEY(record_id) REFERENCES records(id)" +
	");" +
	"CREATE TABLE IF NOT EXISTS audits (" +
	" id INTEGER PRIMARY KEY AUTOINCREMENT," +
	" actor TEXT NOT NULL," +
	" action TEXT NOT NULL," +
	" detail TEXT NOT NULL," +
	" created_at TEXT NOT NULL" +
	");" +
	"CREATE TABLE IF NOT EXISTS settings (" +
	" k TEXT PRIMARY KEY," +
	" v TEXT NOT NULL" +
	");"
