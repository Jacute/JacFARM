import { useState, useEffect } from "react"
import type { Team } from "../../models/models"
import { PAGE_LIMIT } from "../../api/api"
import { deleteTeam } from "../../api/teams"
import { toast } from "react-toastify"

interface Props {
  page: number
  teams: Team[]
  onStats?: (team: Team) => void
  setTeams: (teams: Team[]) => void
}

export const TeamTable = ({ page, teams, onStats, setTeams }: Props) => {
    const [contextMenu, setContextMenu] = useState<{
        x: number
        y: number
        team: Team | null
    } | null>(null)

    useEffect(() => {
        const handleClick = () => setContextMenu(null);
        document.addEventListener("click", handleClick);
        return () => document.removeEventListener("click", handleClick);
    }, []);

    const handleStatistics = (team: Team) => {
        if (onStats) {
        onStats(team)
        }
        setContextMenu(null);
    }

    const handleContextMenu = (event: React.MouseEvent, team: Team) => {
        event.preventDefault();
        setContextMenu(
            contextMenu === null
                ? { x: event.clientX, y: event.clientY, team: team }
                : null,
        );
    };

    const handleDelete = (team: Team) => {
        try {
            deleteTeam(team.id);
        } catch (e: any) {
            toast.error(`Error deleting team: ${e.message}`);
        }
        setTeams(teams.filter(t => t.id !== team.id));
        setContextMenu(null);
    }

    return (
    <div className="table-responsive shadow rounded-3 m-0 p-0 position-relative">
      <table className="table table-fixed table-hover table-striped table-bordered align-middle text-center mb-0">
        <thead className="table-primary">
          <tr>
            <th className="w-10">№</th>
            <th>Name</th>
            <th>IP</th>
          </tr>
        </thead>
        <tbody>
          {teams.map((team, index) => (
            <tr
              key={team.id}
              onContextMenu={(e) => handleContextMenu(e, team)}
              style={{ cursor: "context-menu" }}
            >
              <td className="fw-bold">
                {(page - 1) * PAGE_LIMIT + index + 1}
              </td>
              <td title={team.name}>{team.name}</td>
              <td title={team.ip}>{team.ip}</td>
            </tr>
          ))}
        </tbody>
      </table>

      {contextMenu && (
                <div
                    className="position-fixed bg-white border rounded shadow p-2"
                    style={{
                        top: contextMenu.y,
                        left: contextMenu.x,
                        zIndex: 1000,
                    }}
                >
                    <button className="btn btn-sm btn-secondary w-100 mb-1" onClick={() => {
                        if (contextMenu.team) handleStatistics(contextMenu.team)
                    }}>
                        Статистика
                    </button>
                    <button className="btn btn-sm btn-danger w-100" onClick={() => {
                        if (contextMenu.team) handleDelete(contextMenu.team)
                    }}>
                        Удалить
                    </button>
                </div>
            )}
    </div>
  )
}
