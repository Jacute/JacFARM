interface props {
    team_id: number
    setTeamId: (team_id: number) => void
    exploit_id: string
    setExploitId: (exploit_id: string) => void
}

export const FlagFilter = (props: props) => {
    return (
        <div className="row mb-2">
        <div className="col-3">
          <label className="form-label">Команда</label>
          <select
            className="form-select"
            value={props.team_id}
            onChange={(e) => {
                const value = Number(e.target.value)
                if (!isNaN(value)) {
                    props.setTeamId(value)
                }
            }}
          >
            <option value="">Все</option>
            <option value="1">team1</option>
            <option value="2">team2</option>
            <option value="3">team3</option>
          </select>
        </div>
        <div className="col-3">
          <label className="form-label">Эксплойт</label>
          <select
            className="form-select"
            value={props.exploit_id}
            onChange={(e) => {
                props.setExploitId(e.target.value);
            }}
          >
            <option value="">Все</option>
            <option value="exploit1">exploit1</option>
            <option value="exploit2">exploit2</option>
            <option value="exploit3">exploit3</option>
          </select>
        </div>
      </div>
    )
}