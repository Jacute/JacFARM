import { useEffect, useState } from "react"
import type { ExploitShort, Status, TeamShort } from "../../models/models"
import { getShortTeams } from "../../api/teams"
import { getShortExploits } from "../../api/exploits"
import { getFlagStatuses } from "../../api/flags"
import { Selector } from "../Selector"

interface props {
    team_id: number | null
    setTeamId: (team_id: number | null) => void
    exploit_id: string | null
    setExploitId: (exploit_id: string | null) => void
    status_id: number | null
    setStatusId: (status_id: number | null) => void
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
          <Selector
            label="Team"
            value={props.team_id}
            setValue={props.setTeamId}
            options={teams?.map((team) => ({ label: team.ip, value: team.id }))}
          >
          </Selector>
          <Selector
            label="Exploit"
            value={props.exploit_id}
            setValue={props.setExploitId}
            options={exploits?.map((exploit) => ({ label: exploit.name, value: exploit.id }))}
          >
          </Selector>
          <Selector
            label="Status"
            value={props.status_id}
            setValue={props.setStatusId}
            options={statuses?.map((status) => ({ label: status.name, value: status.id }))}
          >
          </Selector>
        <div className="col-3 d-flex align-items-end">
          <button className="btn btn-primary border border-1" onClick={props.loadFlags}>
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth={2}>
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