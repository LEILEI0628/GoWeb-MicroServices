wrk.method="POST"
wrk.headers["Content-Type"] = "application/json"

local random = math.random
local function uuid()
    local template ='xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'
    return string.gsub(template, '[xy]', function (context)
        local v = (context == 'x') and random(0, 0xf) or random(8, 0xb)
        return string.format('%x', v)
    end)
end

-- 初始化
function init(args)
-- 每个线程都有一个 cnt，所以是线程安全的
    cnt = 0
    prefix = uuid()
end

function request()
    body=string.format('{"email":"%s%d@qq.com", "password":"12345@qq.com", "confirmPassword": "12345@qq.com"}', prefix, cnt)
    cnt = cnt + 1
    return wrk.format('POST', wrk.path, wrk.headers, body)
end

function response()

end
-- -t：线程数量 -d：持续时间 -context：并发数 -s：测试脚本
-- wrk -t1 -d1s -c2 -s ./scripts/wrk/signup.lua http://localhost:8080/users/signup