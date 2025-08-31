import { useEffect, useState } from "react"
import type { ExploitShort, Status, TeamShort } from "../../models/models"
import { getShortTeams } from "../../api/teams"
import { getShortExploits } from "../../api/exploits"
import { getFlagStatuses } from "../../api/flags"

interface props {
    team_id: number
    setTeamId: (team_id: number) => void
    exploit_id: string
    setExploitId: (exploit_id: string) => void
    status_id: number
    setStatusId: (status_id: number) => void
    loadFlags: () => void
}

export const FlagFilter = (props: props) => {
    const [teams, setTeams] = useState<TeamShort[]>([]);
    const [exploits, setExploits] = useState<ExploitShort[]>([]);
    const [statuses, setStatuses] = useState<Status[]>([]);

    useEffect(() => {
        getShortTeams().then(teams => setTeams(teams));
        getShortExploits().then(exploits => setExploits(exploits));
        getFlagStatuses().then(statuses => setStatuses(statuses));
    }, []);

    return (
        <div className="row mb-2">
        <div className="col-3">
          <label className="form-label">Team</label>
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
            <option value="">All</option>
            {teams.map((team) => (
              <option key={team.id} value={team.id}>
                {team.ip}
              </option>
            ))}
          </select>
        </div>
        <div className="col-3">
          <label className="form-label">Exploit</label>
          <select
            className="form-select"
            value={props.exploit_id}
            onChange={(e) => {
                props.setExploitId(e.target.value);
            }}
          >
            <option value="">All</option>
            {exploits.map((exploit) => (
              <option key={exploit.id} value={exploit.id}>
                {exploit.name}
              </option>
            ))}
          </select>
        </div>
        <div className="col-3">
          <label className="form-label">Status</label>
          <select
            className="form-select"
            value={props.status_id}
            onChange={(e) => {
                const value = Number(e.target.value)
                if (!isNaN(value)) {
                    props.setStatusId(value)
                }
            }}
          >
            <option value="">All</option>
            {statuses.map((status) => (
              <option key={status.id} value={status.id}>
                {status.name}
              </option>
            ))}
          </select>
        </div>
        <div className="col-3 d-flex align-items-end">
          <button className="btn btn-primary border border-1" onClick={props.loadFlags}>
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M23 4v6h-6"></path>
              <path d="M1 20v-6h6"></path>
              <path d="M2.51 9a9 9 0 0 1 14.85-3.36L23 10"></path>
              <path d="M21.49 15a9 9 0 0 1-14.85 3.36L1 14"></path>
            </svg>
          </button>
        </div>
      </div>
    )
}