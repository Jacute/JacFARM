interface props {
    title: string;
    ModalBody: React.ReactElement<any>;
}

export const Modal = (props: props) => {
    return (
        <div className="modal fade" id="modal" tabIndex={-1} aria-labelledby="exampleModalLabel" aria-hidden="true">
            <div className="modal-dialog">
                <div className="modal-content">
                <div className="modal-header">
                    <h1 className="modal-title fs-5" id="exampleModalLabel">{props.title}</h1>
                    <button type="button" className="btn-close" data-bs-dismiss="modal" aria-label="Закрыть"></button>
                </div>
                {props.ModalBody && <>{props.ModalBody}</>}
                </div>
            </div>
        </div>
    )
}