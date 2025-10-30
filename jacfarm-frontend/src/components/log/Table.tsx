import { useEffect } from 'react'
import type { Log } from '../../models/models'
import { PAGE_LIMIT } from '../../api/api'


interface props {
    page: number
    module_id: number | null
    exploit_id: string | null
    logLevelId: number | null
    setCount: (count: number) => void
    logs: Log[]
    loadLogs: () => void
}

export const LogTable = (props: props) => {
    useEffect(() => {
        props.loadLogs();
    }, [props.module_id, props.page, props.exploit_id, props.logLevelId, props.exploit_id]);

    return (
        <div className="table-responsive shadow rounded-3 m-0 p-0">
            <table className="table table-fixed table-hover table-striped table-bordered align-middle text-center mb-0">
                <thead className="table-primary">
                    <tr>
                        <th className="w-10">№</th>
                        <th>Module</th>
                        <th>Operation</th>
                        <th>Log Level</th>
                        <th>Value</th>
                        <th>Exploit</th>
                        <th>Attributes</th>
                        <th>Created At</th>
                    </tr>
                </thead>
                <tbody>
                    {props.logs?.map(log => (
                        <tr key={log.id}>
                            <td className="fw-bold">{(props.page - 1) * PAGE_LIMIT + props.logs.indexOf(log) + 1}</td>
                            <td title={log.module}>{log.module}</td>
                            <td title={log.operation}>{log.operation}</td>
                            <td title={log.log_level}>
                                {
                                    <span className={(() => {
                                        switch (log.level) {
                                        case "DEBUG":
                                            return "badge bg-success";
                                        case "ERROR":
                                            return "badge bg-danger";
                                        case "WARNING":
                                            return "badge bg-warning"
                                        default:
                                            return "badge bg-info";
                                        }
                                    })()}>
                                        {log.log_level}
                                    </span>
                                }
                            </td>
                            <td title={log.value}>{log.value}</td>
                            <td title={log.exploit}>{log.exploit}</td>
                            <td title={log.attrs}>{log.attrs}</td>
                            <td className="text-muted small" title={log.created_at}>{log.created_at}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    )
}
