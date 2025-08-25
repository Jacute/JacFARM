import { ToastContainer } from 'react-toastify';
import { Menu } from './components/Menu'
import { ExploitsPage } from './pages/Exploits';
import { FlagsPage } from './pages/Flags';
import './styles/index.css'
import { useState } from 'react'

function App() {
  const [page, setPage] = useState<string>("flags");

  return (
    <div className="container-fluid vh-100">
      <div className="row h-100">
        <Menu setPage={setPage}/>
        <div className="col-10">
          {page === "flags" && <FlagsPage />}
          {page === "exploits" && <ExploitsPage />}
          {page === "teams" && <div>Teams</div>}
          {page === "config" && <div>Config</div>}
        </div>
      </div>
      <ToastContainer />
    </div>
  )
}

export default App
