import { useEffect, useState } from 'react'
import type { Flag } from '../../models/models'
import { getFlags, pageLimit } from '../../api/flags'


interface props {
    page: number
    team_id: number
    exploit_id: string
}

export const FlagTable = (props: props) => {
    const [flags, setFlags] = useState<Array<Flag>>([]);

    useEffect(() => {
        getFlags(props.page).then(flags => setFlags(flags));
    }, []);

    return (
        <div className="table-responsive shadow rounded-3 m-0 p-0">
            <table className="table table-hover table-striped table-bordered align-middle text-center mb-0">
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
                    {flags.map(flag => (
                        <tr key={flag.id}>
                            <td className="fw-bold">{(props.page - 1) * pageLimit + flags.indexOf(flag) + 1}</td>
                            <td>{flag.value}</td>
                            <td>
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
                            <td>
                                {flag.exploit_name}
                            </td>
                            <td>
                                {flag.victim_team}
                            </td>
                            <td>{flag.message_from_server}</td>
                            <td className="text-muted small">{flag.created_at}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    )
}
