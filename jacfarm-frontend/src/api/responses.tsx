import type { ExploitShort, Flag, TeamShort } from "../models/models";
import type { BaseResponse } from "./api";

export interface GetFlagsResponse extends BaseResponse {
    flags: Array<Flag>
    count: number
}

export interface GetExploitsShortResponse extends BaseResponse {
    exploits: Array<ExploitShort>
}

export interface GetExploitsResponse extends BaseResponse {
    exploits: Array<ExploitShort>
}

export interface GetShortTeamsResponse extends BaseResponse {
    teams: Array<TeamShort>
}