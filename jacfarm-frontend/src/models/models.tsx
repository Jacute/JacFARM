export interface Flag {
    id: number
    value: string
    status: string
    exploit_name: string
    victim_team: string
    message_from_server: string
    created_at: string
}

export interface Status {
    id: number
    name: string
}

export interface ExploitShort {
    id: string
    name: string
}

export interface Exploit {
    id: string
    name: string
    is_running: boolean
    is_running_on_farm: boolean
    executable_path: string
    requirements_path?: string
    type: string
}

export interface Team {
    id: number
    name: string
    ip: string
}

export interface TeamShort {
    id: number
    ip: string
}

export interface Config {
    id: number
    name: string
    value: string
}

export interface Log {
    id: number
    module: string
    operation: string
    level: string
    value: string
    exploit: string
    attrs: string
    created_at: string
    log_level: string
}

export interface Module {
    id: number
    name: string
}

export interface LogLevel {
    id: number
    name: string
}