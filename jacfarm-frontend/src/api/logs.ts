import { API_URI, PAGE_LIMIT, STATUS_OK } from "./api";
import type { GetLogLevelsResponse, GetLogsResponse, GetModulesResponse } from "./responses";


export const getLogs = async (page: number, moduleId: number | null, exploitId: string | null, logLevelId: number | null) => {
    if (moduleId == null) moduleId = 0;
    if (exploitId == null) exploitId = "";
    if (logLevelId == null) logLevelId = 0;

    const res = await fetch(`${API_URI}/api/v1/logs?page=${page}&limit=${PAGE_LIMIT}&module_id=${moduleId}&exploit_id=${exploitId}&log_level_id=${logLevelId}`, {credentials: "include"});
    const data: GetLogsResponse = await res.json();

    if (data.status != STATUS_OK) {
        throw Error(data.error);
    }

    return {
        logs: data.logs,
        count: data.count
    }
}

export const getModules = async () => {
    const res = await fetch(`${API_URI}/api/v1/logs/modules`, {credentials: "include"});
    const data: GetModulesResponse = await res.json();

    if (data.status != STATUS_OK) {
        throw Error(data.error);
    }

    return {
        modules: data.modules,
        count: data.count
    }
}

export const getLogLevels = async () => {
    const res = await fetch(`${API_URI}/api/v1/logs/levels`, {credentials: "include"});
    const data: GetLogLevelsResponse = await res.json();

    if (data.status != STATUS_OK) {
        throw Error(data.error);
    }

    return {
        log_levels: data.log_levels,
        count: data.count
    }
}