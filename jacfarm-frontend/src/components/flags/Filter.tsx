import { useEffect, useState } from "react"
import type { ExploitShort, TeamShort } from "../../models/models"
import { getShortTeams } from "../../api/teams"
import { getShortExploits } from "../../api/exploits"

interface props {
    team_id: number
    setTeamId: (team_id: number) => void
    exploit_id: string
    setExploitId: (exploit_id: string) => void
}

export const FlagFilter = (props: props) => {
    const [teams, setTeams] = useState<TeamShort[]>([]);
    const [exploits, setExploits] = useState<ExploitShort[]>([]);

    useEffect(() => {
        getShortTeams().then(teams => setTeams(teams));
        getShortExploits().then(exploits => setExploits(exploits));
    }, []);

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
            {teams.map((team) => (
              <option key={team.id} value={team.id}>
                {team.ip}
              </option>
            ))}
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
            {exploits.map((exploit) => (
              <option key={exploit.id} value={exploit.id}>
                {exploit.name}
              </option>
            ))}
          </select>
        </div>
      </div>
    )
}