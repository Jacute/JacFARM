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
    id: number
    name: string
}

export interface Exploit {
    id: number
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