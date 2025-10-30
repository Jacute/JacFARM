import { useCallback, useState } from "react";
import type { Log } from "../models/models";
import { LogFilter } from "../components/log/Filter";
import { getLogs } from "../api/logs";
import { toast } from "react-toastify";
import { Paginator } from "../components/Paginator";
import { Spinner } from "../components/common/Spinner";
import { LogTable } from "../components/log/Table";

export const logsPage = "logs"


export const LogsPage = () => {
    const [page, setPage] = useState<number>(1);
    const [count, setCount] = useState<number>(0);
    const [logs, setLogs] = useState<Array<Log>>([]);
    const [moduleId, setModuleId] = useState<number | null>(null);
    const [exploitId, setExploitId] = useState<string | null>(null);
    const [logLevelId, setLogLevelId] = useState<number | null>(null);
    const [isLoading, setIsLoading] = useState<boolean>(false);

    const loadLogs = useCallback(async () => {
        setIsLoading(true);
        try {
            const { logs, count } = await getLogs(page, moduleId, exploitId, logLevelId);
            setCount(count);
            setLogs(logs);
        } catch (error: any) {
            toast.error(`Error loading exploits: ${error.message}`);
        } finally {
            setIsLoading(false);
        }
    }, [page, moduleId, exploitId, logLevelId, setCount]);

    return (
        <>
            <div className="row h-10 bg-primary border-bottom border-1 border-secondary">
                <div className="col-11"></div>
                <div className="col-1 d-flex align-items-center justify-content-center">
                    <button className='btn btn-white w-50 border border-1' data-bs-toggle="modal" data-bs-target="#modal">
                    ➕
                    </button>
                </div>
            </div>
            <div className="row h-90">
                <div className="d-flex flex-column h-100">
                    <LogFilter
                        moduleId={moduleId}
                        setModuleId={setModuleId}
                        setExploitId={setExploitId}
                        exploitId={exploitId}
                        loadLogs={loadLogs}
                        logLevelId={logLevelId}
                        setLogLevelId={setLogLevelId}
                    >
                    </LogFilter>
                    <div className="row flex-grow-1 overflow-auto">
                        {isLoading && <Spinner></Spinner>}

                        <LogTable
                            logs={logs}
                            setCount={setCount}
                            page={page}
                            module_id={moduleId}
                            exploit_id={exploitId}
                            logLevelId={logLevelId}
                            loadLogs={loadLogs}
                        >
                        </LogTable>
                    </div>

                    <div className="row mt-auto mx-auto">
                        <div className="col">
                            <Paginator page={page} setPage={setPage} count={count}/>
                        </div>
                    </div>
                </div>
            </div>
        </>
    )
}