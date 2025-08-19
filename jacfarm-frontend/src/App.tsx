import { Menu } from './components/Menu'
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
          {page === "exploits" && <div>Exploits</div>}
          {page === "teams" && <div>Teams</div>}
          {page === "config" && <div>Config</div>}
        </div>
      </div>
    </div>
  )
}

export default App
