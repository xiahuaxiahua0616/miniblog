package auth

import (
	"time"

	"golang.org/x/crypto/bcrypt"

	casbin "github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	adapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

const (
	// 默认的 Casbin 访问控制模型.
	defaultAclModel = `[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act, eft

[role_definition]
g = _, _

[policy_effect]
e = !some(where (p.eft == deny))

[matchers]
m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && r.act == p.act`
)

// Encrypt 使用 bcrypt 加密纯文本.
func Encrypt(source string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(source), bcrypt.DefaultCost)

	return string(hashedBytes), err
}

// Compare 比较密文和明文是否相同.
func Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

type Authz struct {
	*casbin.SyncedEnforcer
}

type Option func(*authzConfig)

type authzConfig struct {
	aclModel           string        // casbin的模型字符串
	autoLoadPolicyTime time.Duration // 自动加载策略的时间间隔
}

// defaultAuthzConfig 返回一个默认的配置.
func defaultAuthzConfig() *authzConfig {
	return &authzConfig{
		// 默认使用内置的ACL模型
		aclModel:           defaultAclModel,
		autoLoadPolicyTime: 5 * time.Second,
	}
}

// DefaultOptions 提供默认的授权器选项配置.
func DefaultOptions() []Option {
	return []Option{
		// 使用默认的 ACL 模型
		WithAclModel(defaultAclModel),
		// 设置自动加载策略的时间间隔为 10 秒
		WithAutoLoadPolicyTime(10 * time.Second),
	}
}

// WithAclModel 允许通过选项自定义 ACL 模型.
func WithAclModel(model string) Option {
	return func(cfg *authzConfig) {
		cfg.aclModel = model
	}
}

// WithAutoLoadPolicyTime 允许通过选项自定义自动加载策略的时间间隔.
func WithAutoLoadPolicyTime(interval time.Duration) Option {
	return func(cfg *authzConfig) {
		cfg.autoLoadPolicyTime = interval
	}
}

// NewAuthz 创建一个使用 Casbin 完成授权的授权器，通过函数选项模式支持自定义配置.
func NewAuthz(db *gorm.DB, opts ...Option) (*Authz, error) {
	// 初始化默认配置
	cfg := defaultAuthzConfig()

	// 应用所有选项
	for _, opt := range opts {
		opt(cfg)
	}

	// 初始化 Gorm 适配器并用于 Casbin 授权器
	adapter, err := adapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err // 返回错误
	}

	// 从配置中加载 Casbin 模型
	m, _ := model.NewModelFromString(cfg.aclModel)

	// 初始化授权器
	enforcer, err := casbin.NewSyncedEnforcer(m, adapter)
	if err != nil {
		return nil, err // 返回错误
	}

	// 从数据库加载策略
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err // 返回错误
	}

	// 启动自动加载策略，使用配置的时间间隔
	enforcer.StartAutoLoadPolicy(cfg.autoLoadPolicyTime)

	// 返回新的授权器实例
	return &Authz{enforcer}, nil
}

// Authorize 用于进行授权.
func (a *Authz) Authorize(sub, obj, act string) (bool, error) {
	// 调用 Enforce 方法进行授权检查
	return a.Enforce(sub, obj, act)
}
