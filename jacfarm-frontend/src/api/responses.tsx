import type { ExploitShort, Flag } from "../models/models";
import type { BaseResponse } from "./api";

export interface GetFlagsResponse extends BaseResponse {
    flags: Array<Flag>
}

export interface GetExploitsShortResponse extends BaseResponse {
    exploits: Array<ExploitShort>
}

export interface GetExploitsResponse extends BaseResponse {
    exploits: Array<ExploitShort>
}