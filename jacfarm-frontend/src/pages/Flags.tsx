import { Paginator } from "../components/Paginator"
import { FlagTable } from "../components/flags/Table"
import { FlagFilter } from "../components/flags/Filter"
import { useState } from "react"
import { AddModal } from "../components/flags/AddModal"

export const flagsPage = "flags"

export const FlagsPage = () => {
  const [team_id, setTeamId] = useState<number>(0);
  const [exploit_id, setExploitId] = useState<string>("");
  const [status_id, setStatusId] = useState<number>(0);
  const [page, setPage] = useState<number>(1);
  const [count, setCount] = useState<number>(0);

  return (
    <>
        <div className="row h-10 bg-primary border-bottom border-1 border-secondary">
            <div className="col-11"></div>
            <div className="col-1 d-flex align-items-center justify-content-center">
                <button className='btn btn-white w-50 border border-1' data-bs-toggle="modal" data-bs-target="#exampleModal">
                ➕
                </button>
            </div>
        </div>
        <div className="row h-90">
            <div className="d-flex flex-column h-100">
                <FlagFilter
                    team_id={team_id}
                    setTeamId={setTeamId}
                    exploit_id={exploit_id}
                    setExploitId={setExploitId}
                    status_id={status_id}
                    setStatusId={setStatusId}
                />

                <div className="row flex-grow-1 overflow-auto">
                    <FlagTable page={page} team_id={team_id} exploit_id={exploit_id} status_id={status_id} setCount={setCount} />
                </div>

                <div className="row mt-auto mx-auto">
                    <div className="col">
                        <Paginator page={page} setPage={setPage} count={count}/>
                    </div>
                </div>
            </div>
        </div>
        <AddModal/>
    </>
  )
}
