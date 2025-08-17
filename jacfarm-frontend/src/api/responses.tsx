import type { Flag } from "../models/models";
import type { BaseResponse } from "./api";

export interface GetFlagsResponse extends BaseResponse {
    flags: Array<Flag>
}