pwd && (
mockgen -source=./webook/internal/web/jwt/types.go \
  -package=jwtmocks -destination=./webook/internal/web/jwt/mocks/handler_mock.go

mockgen -source=./webook/internal/service/user.go \
  -package=svcmocks -destination=./webook/internal/service/mocks/user_mock.go
mockgen -source=./webook/internal/service/code.go \
  -package=svcmocks -destination=./webook/internal/service/mocks/code_mock.go
mockgen -source=./webook/internal/service/article.go \
  -package=svcmocks -destination=./webook/internal/service/mocks/article_mock.go
mockgen -source=./webook/internal/service/sms/types.go \
  -package=smsmocks -destination=./webook/internal/service/sms/mocks/svc_mock.go
mockgen -source=./webook/internal/service/oauth2/wechat/wechat.go \
  -package=wechatmocks -destination=./webook/internal/service/oauth2/wechat/mocks/svc_mock.go

mockgen -source=./webook/internal/repository/code.go \
  -package=repomocks -destination=./webook/internal/repository/mocks/code_mock.go
mockgen -source=./webook/internal/repository/user.go \
  -package=repomocks -destination=./webook/internal/repository/mocks/user_mock.go
mockgen -source=./webook/internal/repository/article.go \
  -package=repomocks -destination=./webook/internal/repository/mocks/article_mock.go

mockgen -source=./webook/internal/repository/dao/user.go \
  -package=daomocks -destination=./webook/internal/repository/dao/mocks/user_mock.go
mockgen -source=./webook/internal/repository/dao/article.go \
  -package=daomocks -destination=./webook/internal/repository/dao/mocks/article_mock.go
mockgen -source=./webook/internal/repository/cache/user.go \
  -package=cachemocks -destination=./webook/internal/repository/cache/mocks/user_mock.go

mockgen -source=./webook/pkg/limiter/types.go \
  -package=limitmocks -destination=./webook/pkg/limiter/mocks/limiter_mock.go

mockgen -package=redismocks -destination=./webook/internal/repository/cache/redismocks/cmd_mock.go github.com/redis/go-redis/v9 Cmdable
)
