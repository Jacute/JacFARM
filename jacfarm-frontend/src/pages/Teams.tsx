import { useCallback, useEffect, useState } from "react"
import { Paginator } from "../components/Paginator"
import { Modal } from "../components/Modal"
import { AddTeamModal } from "../components/teams/AddModal"
import { TeamTable } from "../components/teams/Table"
import type { Team } from "../models/models"
import { getTeams } from "../api/teams"
import { toast } from "react-toastify"

export const teamsPage = "teams"

export const TeamsPage = () => {
  const [search, setSearch] = useState<string>("")
  const [debouncedSearch, setDebouncedSearch] = useState<string>("")
  const [page, setPage] = useState<number>(1)
  const [count, setCount] = useState<number>(0)
  const [teams, setTeams] = useState<Array<Team>>([])
  const [isLoading, setIsLoading] = useState<boolean>(false)

  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedSearch(search)
      setPage(1)
    }, 300)

    return () => clearTimeout(handler)
  }, [search])

  const loadTeams = useCallback(async () => {
    setIsLoading(true)
    try {
      const { teams, count } = await getTeams(page, debouncedSearch)
      setCount(count)
      setTeams(teams)
    } catch (error: any) {
      toast.error(`Error loading teams: ${error.message}`)
    } finally {
      setIsLoading(false)
    }
  }, [page, debouncedSearch])

  useEffect(() => {
    loadTeams()
  }, [loadTeams])

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
            <div className="col-10 d-flex align-items-center my-3">
                <input
                type="text"
                placeholder="Search teams..."
                value={search}
                onChange={e => setSearch(e.target.value)}
                className="form-control w-50"
                />
                <button
                className="btn btn-primary border border-1 ms-2"
                onClick={() => {
                    setDebouncedSearch(search)
                    setPage(1)
                }}
                >
                🔍
                </button>
            </div>

            <div className="row flex-grow-1 overflow-auto position-relative">
                {isLoading && (
                <div
                    className="position-absolute top-0 start-0 w-100 h-100 d-flex align-items-center justify-content-center"
                    style={{
                    backgroundColor: "rgba(255, 255, 255, 0.8)",
                    zIndex: 1000,
                    }}
                >
                    <div className="spinner-border text-primary" role="status">
                    <span className="visually-hidden">Loading...</span>
                    </div>
                </div>
                )}

                <TeamTable page={page} teams={teams} setTeams={setTeams} />
            </div>

            <div className="row mt-auto mx-auto">
                <div className="col">
                <Paginator page={page} setPage={setPage} count={count} />
                </div>
            </div>
            </div>
        </div>

        <Modal
            title="Добавить команду"
            ModalBody={<AddTeamModal loadTeams={loadTeams} />}
        />
    </>
  )
}
