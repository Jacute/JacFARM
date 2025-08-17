export const API_URI = import.meta.env.VITE_API_URI;

export const STATUS_OK = "OK";
export const STATUS_ERROR = "Error";

export interface BaseResponse {
    status: string;
    error?: string;
};