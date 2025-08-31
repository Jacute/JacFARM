import { ToastContainer } from 'react-toastify';
import { Menu } from './components/Menu'
import { ExploitsPage } from './pages/Exploits';
import { FlagsPage } from './pages/Flags';
import './styles/index.css'
import { useState } from 'react'
import { TeamsPage } from './pages/Teams';
import { ConfigPage } from './pages/Config';

function App() {
  const [page, setPage] = useState<string>("flags");

  return (
    <div className="container-fluid vh-100">
      <div className="row h-100">
        <Menu setPage={setPage}/>
        <div className="col-10">
          {page === "flags" && <FlagsPage />}
          {page === "exploits" && <ExploitsPage />}
          {page === "teams" && <TeamsPage />}
          {page === "config" && <ConfigPage />}
        </div>
      </div>
      <ToastContainer />
    </div>
  )
}

export default App
