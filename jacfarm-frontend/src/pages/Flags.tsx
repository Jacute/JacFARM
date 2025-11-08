import { Paginator } from "../components/Paginator"
import { FlagTable } from "../components/flags/Table"
import { FlagFilter } from "../components/flags/Filter"
import { useCallback, useState } from "react"
import { AddModal } from "../components/flags/AddModal"
import { Modal } from "../components/Modal"
import type { Flag } from "../models/models"
import { getFlags } from "../api/flags"
import { toast } from "react-toastify"
import { Spinner } from "../components/common/Spinner"


export const FlagsPage = () => {
  const [team_id, setTeamId] = useState<number | null>(null);
  const [exploit_id, setExploitId] = useState<string | null>(null);
  const [status_id, setStatusId] = useState<number | null>(null);
  const [page, setPage] = useState<number>(1);
  const [count, setCount] = useState<number>(0);
  const [flags, setFlags] = useState<Array<Flag>>([]);
  const [isLoading, setIsLoading] = useState<boolean>(false); // Состояние загрузки

  const loadFlags = useCallback(async () => {
    setIsLoading(true);
    try {
        const { flags, count } = await getFlags(page, team_id, exploit_id, status_id);
        setCount(count);
        setFlags(flags);
    } catch (error: any) {
        toast.error(`Error loading flags: ${error.message}`);
    } finally {
        setIsLoading(false);
    }
  }, [page, team_id, exploit_id, status_id, setCount]);

  return (
    <>
        <div className="row h-10 bg-primary border-bottom border-1 border-secondary">
            <div className="col-11"></div>
            <div className="col-1 d-flex align-items-center justify-content-center">
                <button className='btn btn-white w-50 border border-1' data-bs-toggle="modal" data-bs-target="#modal">
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
                    loadFlags={loadFlags}
                />

                <div className="row flex-grow-1 overflow-auto position-relative">
                    {isLoading && <Spinner></Spinner>}
                    
                    <FlagTable 
                        page={page} 
                        team_id={team_id} 
                        exploit_id={exploit_id} 
                        status_id={status_id} 
                        setCount={setCount} 
                        flags={flags} 
                        loadFlags={loadFlags} 
                    />
                </div>

                <div className="row mt-auto mx-auto">
                    <div className="col">
                        <Paginator page={page} setPage={setPage} count={count}/>
                    </div>
                </div>
            </div>
        </div>
        <Modal title="Добавить флаг" ModalBody={<AddModal loadFlags={loadFlags}/>}/>
    </>
  )
}