import { useState } from "react"
import type { Config } from "../../models/models"
import { PAGE_LIMIT } from "../../api/api"
import { toast } from "react-toastify"
import { updateConfig } from "../../api/config"

interface Props {
  page: number
  configs: Config[]
  setConfigs: (configs: Config[]) => void
}

export const ConfigTable = ({ page, configs, setConfigs }: Props) => {
  const [editing, setEditing] = useState<{ id: number; field: "key" | "value" } | null>(null)
  const [editValue, setEditValue] = useState("")

  const handleDoubleClick = (id: number, field: "key" | "value", currentValue: string) => {
    setEditing({ id, field })
    setEditValue(currentValue)
  }

  const handleBlur = async () => {
    if (!editing) return

    try {
        await updateConfig(editing.id, editValue)
        setConfigs(configs.map(cfg => (cfg.id === editing.id ? { ...cfg, value: editValue } : cfg)))
    } catch (e: any) {
        toast.error(`Ошибка при обновлении: ${e.message}`)
    }
    setEditing(null)
    setEditValue("")
  }

  return (
    <div className="table-responsive shadow rounded-3 m-0 p-0 position-relative">
      <table className="table table-hover table-striped table-bordered align-middle text-center mb-0">
        <thead className="table-primary">
          <tr>
            <th className="w-10">№</th>
            <th>Key</th>
            <th>Value</th>
          </tr>
        </thead>
        <tbody>
          {configs.map((cfg, index) => (
            <tr key={cfg.id}>
              <td className="fw-bold">{(page - 1) * PAGE_LIMIT + index + 1}</td>
              <td
                onDoubleClick={() => handleDoubleClick(cfg.id, "key", cfg.name)}
                style={{ cursor: "pointer" }}
              >
                {cfg.name}
              </td>
              <td
                onDoubleClick={() => handleDoubleClick(cfg.id, "value", cfg.value)}
                style={{ cursor: "pointer" }}
              >
                {editing?.id === cfg.id && editing.field === "value" ? (
                  <input
                    autoFocus
                    type="text"
                    className="w-100"
                    value={editValue}
                    onChange={(e) => setEditValue(e.target.value)}
                    onBlur={handleBlur}
                    onKeyDown={(e) => e.key === "Enter" && handleBlur()}
                  />
                ) : (
                  cfg.value
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
