package defs

const (
	CapFileName    = "workerd.capnp"
	WorkerInfoPath = "workers"
	WorkerCodePath = "src"
	DBTypeSqlite   = "sqlite"
	DBTypeMysql    = "mysql"
	DBTypePostgres = "pgsql"

	DefaultHostName     = "localhost"
	DefaultNodeName     = "default"
	DefaultExternalPath = "/"
	DefaultEntry        = "entry.js"
	DefaultCode         = `export default {
  async fetch(req, env) {
    try {
		let resp = new Response("worker: " + req.url + " is online! -- " + new Date())
		return resp
	} catch(e) {
		return new Response(e.stack, { status: 500 })
	}
  }
};`

	DefaultTemplate = `using Workerd = import "/workerd/workerd.capnp";

const config :Workerd.Config = (
  services = [
    ( name = "{{.Worker.UID}}", worker = .v{{.Worker.UID}}Worker),
	( 
		name = "internet", 
		network = (allow = ["public", "private"],
	 	tlsOptions = (trustBrowserCas = true)
	)),
	{{.ServiceText}}
  ],

  sockets = [
    (
      name = "{{.Worker.UID}}",
      address = "{{.Worker.HostName}}:{{.Worker.Port}}",
      http=(),
      service="{{.Worker.UID}}"
    ),
  ],

  extensions = [{{.ExtensionsText}}],
);

const v{{.Worker.UID}}Worker :Workerd.Worker = (
  modules = [
    (name = "{{.Worker.Entry}}", esModule = embed "src/{{.Worker.Entry}}"),
  ],
  compatibilityDate = "2025-05-08",
  bindings = [{{.BindingsText}}],
  compatibilityFlags = [{{.FlagsText}}],
);

`
)

const (
	KeyNodeName    = "node_name"
	KeyNodeSecret  = "node_secret"
	KeyNodeProto   = "node_proto"
	KeyWorkerProto = "worker_proto"
)

const (
	HeaderNodeName   = "x-node-name"
	HeaderNodeSecret = "x-secret"
	HeaderHost       = "Host"
)

const (
	CodeInvalidRequest = 400
	CodeUnAuthorized   = 401
	CodeNotFound       = 404
	CodeInternalError  = 500
	CodeSuccess        = 200
)

const (
	EventSyncWorkers  = "sync-workers"
	EventAddWorker    = "add-worker"
	EventDeleteWorker = "delete-worker"
	EventFlushWorker  = "flush-worker"
)
