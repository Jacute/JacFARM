import { useEffect } from 'react'
import type { Flag } from '../../models/models'
import { PAGE_LIMIT } from '../../api/api'


interface props {
    page: number
    team_id: number | null
    exploit_id: string | null
    status_id: number | null
    setCount: (count: number) => void
    flags: Flag[]
    loadFlags: () => void
}

export const FlagTable = (props: props) => {
    useEffect(() => {
        props.loadFlags()
    }, [props.team_id, props.page, props.exploit_id, props.status_id]);

    return (
        <div className="table-responsive shadow rounded-3 m-0 p-0">
            <table className="table table-fixed table-hover table-striped table-bordered align-middle text-center mb-0">
                <thead className="table-primary">
                    <tr>
                        <th className="w-10">№</th>
                        <th>Значение</th>
                        <th>Статус</th>
                        <th>Эксплойт</th>
                        <th>Команда</th>
                        <th>Ответ от сервера</th>
                        <th>Дата создания</th>
                    </tr>
                </thead>
                <tbody>
                    {props.flags.map(flag => (
                        <tr key={flag.id}>
                            <td className="fw-bold">{(props.page - 1) * PAGE_LIMIT + props.flags.indexOf(flag) + 1}</td>
                            <td title={flag.value}>{flag.value}</td>
                            <td title={flag.status}>
                                <span className={(() => {
                                    switch (flag.status) {
                                    case "SUCCESS":
                                        return "badge bg-success";
                                    case "OLD":
                                        return "badge bg-danger";
                                    default:
                                        return "badge bg-info";
                                    }
                                })()}>
                                    {flag.status}
                                </span>
                            </td>
                            <td title={flag.exploit_name}>
                                {flag.exploit_name}
                            </td>
                            <td title={flag.victim_team}>
                                {flag.victim_team}
                            </td>
                            <td title={flag.message_from_server}>{flag.message_from_server}</td>
                            <td className="text-muted small" title={flag.created_at}>{flag.created_at}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    )
}
