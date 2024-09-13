namespace go user

struct User {
    1: string id,
    2: string username,
    3: string email
    4: map<string, UserStatus> status
}

struct UserStatus {
    1: string deviceID,
    2: string serverAddress
}

service UserService {
    // 注册接口，传入用户名、密码和邮箱，返回成功或失败
    bool Register(1: string username, 2: string password, 3: string email)

    // 登录接口，传入用户名和密码，返回 JWT 令牌
    string Login(1: string username, 2: string password)

    // 验证 JWT 令牌，返回用户信息
    User ValidateToken(1: string token)

    // 获取用户信息，传入用户 ID (UUID)，返回用户信息
    User GetUserInfo(1: string id)

    // 设置用户设备在线，传入用户 ID (UUID)、设备 ID 和服务器地址
    bool SetUserOnline(1: string id, 2: string deviceID, 3: string serverAddress)

    // 设置用户设备离线，传入用户 ID (UUID) 和设备 ID
    bool SetUserOffline(1: string id, 2: string deviceID)

    // 获取用户的所有在线设备及其连接的服务器
    map<string, UserStatus> GetUserDevices(1: string id)
}
