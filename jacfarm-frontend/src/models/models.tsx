export interface Flag {
    id: number
    value: string
    status: string
    exploit_name: string
    victim_team: string
    message_from_server: string
    created_at: string
}

export interface ExploitShort {
    id: number
    name: string
}

export interface Exploit {
    id: number
    name: string
    is_running: string
    is_running_on_farm: string
    executable_path: string
    requirements_path?: string
    type: string
}