package conf

import (
	"fmt"
	"vvorker/utils/secret"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type AppConfig struct {
	MasterEndpoint string `env:"MASTER_ENDPOINT" env-default:"http://127.0.0.1:8888"` // needed for agent，agent需要通过该url来注册节点
	WorkerPort     int    `env:"WORKER_PORT" env-default:"8080"`                      // 【主节点公开】提供worker服务，如 xxx.example.com:8080
	APIPort        int    `env:"API_PORT" env-default:"8888"`                         // 【主节点公开】提供控制台服务和控制api服务

	TunnelHost      string `env:"TUNNEL_HOST" env-default:"127.0.0.1"`   // for master usually 127.0.0.1, for agent usually master public ip
	TunnelEntryPort int    `env:"TUNNEL_ENTRY_PORT" env-default:"10080"` // 【主节点内部】提供http服务，主节点用这个端口向其他节点发送请求。在主节点提供服务，子节点无法向其发送请求。
	TunnelAPIPort   int    `env:"TUNNEL_API_PORT" env-default:"18080"`   // 【主节点公开】TUNNEL_ENTRY_PORT的frp端口，子节点通过这个端口配置每个worker的转发

	WorkerURLSuffix string `env:"WORKER_URL_SUFFIX"`         // master required, e.g. .example.com. for worker show and route
	Scheme          string `env:"SCHEME" env-default:"http"` // http, https. for public frontend show
	NodeName        string `env:"NODE_NAME" env-default:"default"`
	AgentSecret     string `env:"AGENT_SECRET"` //	required, e.g. 123123123

	DBPath         string `env:"DB_PATH" env-default:"/workerd/db.sqlite"`
	WorkerdDir     string `env:"WORKERD_DIR" env-default:"/workerd"`
	DBType         string `env:"DB_TYPE" env-default:"sqlite"`
	WorkerLimit    int    `env:"WORKER_LIMIT" env-default:"10000"`
	WorkerdBinPath string `env:"WORKERD_BIN_PATH" env-default:"/bin/workerd"`

	APIWebBaseURL  string `env:"API_WEB_BASE_URL"`
	ListenAddr     string `env:"LISTEN_ADDR" env-default:"0.0.0.0"`
	CookieName     string `env:"COOKIE_NAME" env-default:"authorization"`
	CookieAge      int    `env:"COOKIE_AGE" env-default:"86400"` // second 86400 = 1 day
	CookieDomain   string `env:"COOKIE_DOMAIN"`                  // required, e.g. example.com
	EnableRegister bool   `env:"ENABLE_REGISTER" env-default:"true"`
	RunMode        string `env:"RUN_MODE" env-default:"master"` // master, agent

	DefaultWorkerHost string `env:"DEFAULT_WORKER_HOST" env-default:"localhost"`
	LitefsPrimaryPort int    `env:"LITEFS_PRIMARY_PORT" env-default:"20202"`
	LitefsBinPath     string `env:"LITEFS_BIN_PATH" env-default:"/usr/local/bin/litefs"`
	LitefsDirPath     string `env:"LITEFS_DIR_PATH" env-default:"/app"`
	LitefsEnabled     bool   `env:"LITEFS_ENABLED" env-default:"false"`
	EnableAutoSync    bool   `env:"ENABLE_AUTO_SYNC" env-default:"false"`
	TunnelUsername    string
	TunnelPassword    string
	TunnelToken       string
	NodeID            string
	WorkerHostMode    string `env:"WORKER_HOST_MODE" env-default:"host"` // host path  // host 模式需要使用域名进行访问，path则url的第一段为服务名（不包含域名后缀

	AdminAPIProxy bool `env:"ADMIN_API_PROXY" env-default:"false"` // 允许admin页面代理api请求，这可能会导致路径冲突，并且WORKER_HOST_MODE必须为path

	ServerRedisHost string `env:"SERVER_REDIS_HOST" env-default:"localhost"`
	ServerRedisPort int    `env:"SERVER_REDIS_PORT" env-default:"6379"`

	ServerMinioHost   string `env:"SERVER_MINIO_HOST" env-default:"localhost"` // localhost时为本地
	ServerMinioPort   int    `env:"SERVER_MINIO_PORT" env-default:"9000"`      // 本地时为9000，远程时为443
	ServerMinioRegion string `env:"SERVER_MINIO_REGION" env-default:"us-east-1"`
	ServerMinioUseSSL bool   `env:"SERVER_MINIO_USE_SSL" env-default:"false"`
	ServerMinioAccess string `env:"SERVER_MINIO_ACCESS" env-default:"minioadmin"`
	ServerMinioSecret string `env:"SERVER_MINIO_SECRET" env-default:"minioadmin"`

	ServerPostgreHost     string `env:"SERVER_POSTGRE_HOST" env-default:"localhost"`
	ServerPostgrePort     int    `env:"SERVER_POSTGRE_PORT" env-default:"5432"`
	ServerPostgrePassword string `env:"SERVER_POSTGRE_PASSWORD" env-default:"postgres"`
	ServerPostgreUser     string `env:"SERVER_POSTGRE_USER" env-default:"postgres"`

	ClientMinioPort   int `env:"CLIENT_MINIO_PORT" env-default:"19000"`
	ClientPostgrePort int `env:"CLIENT_POSTGRE_PORT" env-default:"15432"`
	ClientRedisPort   int `env:"CLIENT_REDIS_PORT" env-default:"16379"`

	LocalTMPPostgrePort int `env:"LOCAL_TMP_POSTGRE_PORT" env-default:"13420"`
	LocalTMPRedisPort   int `env:"LOCAL_TMP_REDIS_PORT" env-default:"13421"`
}

type JwtConfig struct {
	Secret     string `env:"JWT_SECRET" env-default:"secret"`
	ExpireTime int64  `env:"JWT_EXPIRETIME" env-default:"24"` // hour
}

type JwtClaims struct {
	jwt.RegisteredClaims
	UID uint `json:"uid,omitempty"`
}

var (
	AppConfigInstance *AppConfig
	JwtConf           *JwtConfig
	RPCToken          string
)

func init() {
	var err error
	AppConfigInstance = &AppConfig{}
	JwtConf = &JwtConfig{}
	godotenv.Load()

	if err := cleanenv.ReadEnv(AppConfigInstance); err != nil {
		logrus.Panic(err)
	}
	if err := cleanenv.ReadEnv(JwtConf); err != nil {
		logrus.Panic(err)
	}

	RPCToken = secret.MD5(fmt.Sprintf("%s%s", AppConfigInstance.NodeName, AppConfigInstance.AgentSecret))
	AppConfigInstance.TunnelUsername = secret.MD5(AppConfigInstance.AgentSecret +
		AppConfigInstance.WorkerURLSuffix)
	AppConfigInstance.TunnelPassword = secret.MD5(AppConfigInstance.AgentSecret +
		AppConfigInstance.WorkerURLSuffix + AppConfigInstance.TunnelUsername)
	AppConfigInstance.TunnelToken = AppConfigInstance.TunnelPassword

	if err != nil {
		logrus.Panic(err)
	}
}

func IsMaster() bool {
	return AppConfigInstance.RunMode == "master"
}
