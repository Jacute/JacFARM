import { useState } from "react";
import { sendFlag } from "../../api/flags";
import { toast } from "react-toastify";

interface props {
    loadFlags: () => void
}

export const AddModal = (props: props) => {
    const [flag, setFlag] = useState<string>("");

    const handleSubmit = () => {
        try {
            sendFlag(flag);
        } catch (error: any) {
            toast.error(`Error sending flag: ${error.message}`);
            return;
        }
        props.loadFlags();
    }

    return (
        <>
            <div className="modal-body d-flex flex-column">
                <label htmlFor="flag">Flag</label>
                <input
                    type="text"
                    id="flag"
                    className="form-control"
                    value={flag}
                    onChange={(e) => setFlag(e.target.value)}
                />
            </div>
            <div className="modal-footer">
                <button type="button" className="btn btn-secondary" data-bs-dismiss="modal">Close</button>
                <button type="button" className="btn btn-primary" data-bs-dismiss="modal" onClick={handleSubmit}>Save</button>
            </div>
        </>
    )
}