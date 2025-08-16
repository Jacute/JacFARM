import { Menu } from './components/Menu'
import { Table } from './components/Table'
import { Paginator } from './components/Paginator'
import './styles/index.css'

function App() {
  return (
    <div className="container-fluid vh-100">
      <div className="row h-100">
        <Menu></Menu>
        <div className="col-10 d-flex flex-column">
          <div className="row h-10 bg-primary border-bottom border-5 border-secondary">
            Хэдер
          </div>
          <div className="row">
            <div className="row border g-0">
              Фильтр
            </div>
            <div className="row m-0 p-0">
              <Table></Table>
            </div>
          </div>
          <div
              className="row mt-auto mx-auto"
            >
              <div className="col">
                <Paginator></Paginator>
              </div>
            </div>
        </div>
      </div>
    </div>
  )
}

export default App
