import { useState } from "react";
import { sendFlag } from "../../api/flags";
import { toast } from "react-toastify";

export const AddModal = () => {
    const [flag, setFlag] = useState<string>("");

    const handleSubmit = () => {
        try {
            sendFlag(flag);
        } catch (error: any) {
            toast.error(`Error sending flag: ${error.message}`);
        }
    }

    return (
        <>
            <div className="modal-body d-flex flex-column">
                <label htmlFor="flag">Флаг</label>
                <input
                    type="text"
                    id="flag"
                    className="form-control"
                    value={flag}
                    onChange={(e) => setFlag(e.target.value)}
                />
            </div>
            <div className="modal-footer">
                <button type="button" className="btn btn-secondary" data-bs-dismiss="modal">Закрыть</button>
                <button type="button" className="btn btn-primary" data-bs-dismiss="modal" onClick={handleSubmit}>Сохранить</button>
            </div>
        </>
    )
}