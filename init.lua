local log = require('log')

box.cfg{
    listen = 3301,
    wal_mode = 'none',
}

box.schema.user.grant('guest', 'super', nil, nil, { if_not_exists = true })

local last_seq = -1

function reset()
    last_seq = -1
    log.warn('sequence is reset')
end

function test_func(seq)
    local result
    if seq ~= last_seq + 1 then
        result = 'err'
    else
        result = 'ok'
    end

    last_seq = seq
    return result
end
