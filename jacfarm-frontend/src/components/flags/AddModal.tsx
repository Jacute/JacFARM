import { useState } from "react";
import { sendFlag } from "../../api/flags";
import { useRef } from 'react';

export const AddModal = () => {
    const modalRef = useRef(null);
    const [flag, setFlag] = useState<string>("");

    const handleSubmit = () => {
        try {
            sendFlag(flag);
        } catch (error) {
            console.error("Error sending flag:", error);
        }
    }

    return (
        <div className="modal fade" id="exampleModal" tabIndex={-1} aria-labelledby="exampleModalLabel" aria-hidden="true" ref={modalRef}>
        <div className="modal-dialog">
            <div className="modal-content">
            <div className="modal-header">
                <h1 className="modal-title fs-5" id="exampleModalLabel">Добавить флаг</h1>
                <button type="button" className="btn-close" data-bs-dismiss="modal" aria-label="Закрыть"></button>
            </div>
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
            </div>
        </div>
        </div>
    )
}