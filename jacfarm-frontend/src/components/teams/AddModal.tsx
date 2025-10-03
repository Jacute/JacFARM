import { useState } from "react";
import { addTeam } from "../../api/teams";
import { toast } from "react-toastify";

interface props {
    loadTeams: () => void
}

export const AddTeamModal = ({ loadTeams }: props) => {
    const [name, setName] = useState<string>("");
    const [ip, setIp] = useState<string>("");

    const handleSubmit = async () => {
        try {
            await addTeam(name, ip);
            toast.success("Команда успешно добавлена");
            loadTeams();
        } catch (error: any) {
            toast.error(`Ошибка при добавлении команды: ${error.message}`);
            return;
        }
    };

    return (
        <>
            <div className="modal-body d-flex flex-column gap-2">
                <label htmlFor="team-name">Название команды</label>
                <input
                    type="text"
                    id="team-name"
                    className="form-control"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                />

                <label htmlFor="team-ip">IP</label>
                <input
                    id="team-ip"
                    className="form-control"
                    value={ip}
                    onChange={(e) => setIp(e.target.value)}
                />
            </div>

            <div className="modal-footer">
                <button 
                    type="button" 
                    className="btn btn-secondary" 
                    data-bs-dismiss="modal"
                >
                    Закрыть
                </button>
                <button 
                    type="button" 
                    className="btn btn-primary" 
                    data-bs-dismiss="modal" 
                    onClick={handleSubmit}
                >
                    Сохранить
                </button>
            </div>
        </>
    );
};
