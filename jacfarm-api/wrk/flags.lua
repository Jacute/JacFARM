math.randomseed(os.time())

alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
api_key = "TEST_API_KEY"

generate = function ()
    local flag = ""
    for i = 1, 31 do
        local idx = math.random(1, #alphabet)
        flag = flag .. alphabet:sub(idx, idx)
    end
    flag = flag .. "="
    return flag
end

request = function ()
    local flag = generate()
    local body = '{"flag":"' .. flag .. '"}'
    local headers = {
        ["Content-Type"] = "application/json",
        ["Authorization"] = api_key
    }
    return wrk.format("POST", "/api/v1/service/flags", headers, body)
end
