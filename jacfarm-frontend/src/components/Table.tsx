import type { Flag } from '../models/models'


const flags: Array<Flag> = [
    {
        id: 1,
        value: "flag1",
        team: "team1",
        exploit: "exploit1",
        created_by: Date.now().toString()
    },
    {
        id: 2,
        value: "flag2",
        team: "team2",
        exploit: "exploit2",
        created_by: Date.now().toString()
    }
]


export const Table = () => {
    return (
        <div className="table-responsive shadow rounded-3 m-0 p-0">
            <table className="table table-hover table-striped table-bordered align-middle text-center mb-0">
                <thead className="table-primary">
                    <tr>
                        <th className="w-10">№</th>
                        <th>Значение</th>
                        <th>Эксплойт</th>
                        <th>Команда</th>
                        <th>Дата создания</th>
                    </tr>
                </thead>
                <tbody>
                    {flags.map(flag => (
                        <tr key={flag.id}>
                            <td className="fw-bold">{flag.id}</td>
                            <td>{flag.value}</td>
                            <td>
                                <span className="badge bg-info text-dark">
                                    {flag.exploit}
                                </span>
                            </td>
                            <td>
                                <span className="badge bg-secondary">
                                    {flag.team}
                                </span>
                            </td>
                            <td className="text-muted small">{flag.created_by}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    )
}
