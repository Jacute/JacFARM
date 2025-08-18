import { PAGE_LIMIT } from "../api/api"

interface Props {
    page: number
    count: number
    setPage: (page: number) => void
}

export const Paginator = ({ page, count, setPage }: Props) => {
    const hasPrev = page > 1;
    const hasNext = count > page * PAGE_LIMIT;

    return (
        <nav>
            <ul className="pagination">
                {hasPrev && (
                    <>
                        <li className="page-item" key="prev">
                            <button
                                className="page-link"
                                aria-label="Предыдущая"
                                onClick={() => setPage(page - 1)}
                            >
                                &laquo;
                            </button>
                        </li>
                        <li className="page-item" key={page - 1}>
                            <button
                                className="page-link"
                                onClick={() => setPage(page - 1)}
                            >
                                {page - 1}
                            </button>
                        </li>
                    </>
                )}

                <li className="page-item active" key={page}>
                    <button
                        className="page-link"
                        onClick={() => setPage(page)}
                    >
                        {page}
                    </button>
                </li>

                {hasNext && (
                    <>
                        <li className="page-item" key={page + 1}>
                            <button
                                className="page-link"
                                onClick={() => setPage(page + 1)}
                            >
                                {page + 1}
                            </button>
                        </li>
                        <li className="page-item" key="next">
                            <button
                                className="page-link"
                                aria-label="Следующая"
                                onClick={() => setPage(page + 1)}
                            >
                                &raquo;
                            </button>
                        </li>
                    </>
                )}
            </ul>
        </nav>
    )
}
