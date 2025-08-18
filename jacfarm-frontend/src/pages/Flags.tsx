import { Paginator } from "../components/Paginator"
import { FlagTable } from "../components/flags/Table"
import { FlagFilter } from "../components/flags/Filter"
import { useState } from "react"

export const flagsPage = "flags"

export const FlagsPage = () => {
  const [team_id, setTeamId] = useState<number>(0);
  const [exploit_id, setExploitId] = useState<string>("");
  const [page, setPage] = useState<number>(1);
  const [count, setCount] = useState<number>(0);

  return (
    <div className="d-flex flex-column h-100">
        <FlagFilter
            team_id={team_id}
            setTeamId={setTeamId}
            exploit_id={exploit_id}
            setExploitId={setExploitId}
        />

        <div className="row flex-grow-1 overflow-auto">
            <FlagTable page={page} team_id={team_id} exploit_id={exploit_id} setCount={setCount} />
        </div>

        <div className="row mt-auto mx-auto">
            <div className="col">
                <Paginator page={page} setPage={setPage} count={count}/>
            </div>
        </div>
    </div>
  )
}
