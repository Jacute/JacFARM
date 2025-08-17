import { Paginator } from "../components/Paginator"
import { FlagTable } from "../components/flags/Table"
import { FlagFilter } from "../components/flags/Filter"
import { useState } from "react"

export const flagsPage = "flags"

export const FlagsPage = () => {
  const [team_id, setTeamId] = useState<number>(0);
  const [exploit_id, setExploitId] = useState<string>("");

  return (
    <div className="d-flex flex-column h-100">
        <FlagFilter
            team_id={team_id}
            setTeamId={setTeamId}
            exploit_id={exploit_id}
            setExploitId={setExploitId}
        />

        <div className="row flex-grow-1 overflow-auto">
            <FlagTable page={1} team_id={team_id} exploit_id={exploit_id} />
        </div>

        <div className="row mt-auto mx-auto">
            <div className="col">
                <Paginator />
            </div>
        </div>
    </div>
  )
}
