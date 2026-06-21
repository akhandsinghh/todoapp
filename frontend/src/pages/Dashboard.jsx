import React, { useEffect, useMemo, useState } from 'react';
import { Link } from 'react-router-dom'; // <-- IMPORT REQUIRED FOR THE BUTTON
import Navbar from '../components/Navbar';
import GroupSidebar from '../components/GroupSidebar';
import TaskForm from '../components/TaskForm';
import TaskList from '../components/TaskList';
import * as taskApi from '../api/taskApi';
import * as groupApi from '../api/groupApi';

export default function Dashboard() {
  const [groups, setGroups] = useState([]);
  const [tasks, setTasks] = useState([]);
  const [activeGroup, setActiveGroup] = useState(null);
  const [status, setStatus] = useState('');
  const [error, setError] = useState('');
  const [page, setPage] = useState(1);
  const [totalTasks, setTotalTasks] = useState(0);
  const [query, setQuery] = useState({ pageSize: 10, sortBy: 'due_at', sortOrder: 'asc' });
  const [controls, setControls] = useState(query);

  const loadGroups = () => groupApi.listGroups().then(setGroups);
  const loadTasks = () =>
    taskApi
      .listTasks({
        status,
        group_id: typeof activeGroup === 'number' ? activeGroup : undefined,
        ungrouped: activeGroup === 'ungrouped' ? true : undefined,
        page,
        limit: query.pageSize,
        sort_by: query.sortBy,
        sort_order: query.sortOrder,
      })
      .then((data) => {
        const items = Array.isArray(data) ? data : data.items || [];
        setTasks(items);
        setTotalTasks(Array.isArray(data) ? items.length : data.total || 0);
      });

  const taskMatchesFilters = (task) => {
    if (status && task.status !== status) return false;
    if (activeGroup === 'ungrouped' && task.group_id) return false;
    if (typeof activeGroup === 'number' && task.group_id !== activeGroup) return false;
    return true;
  };

  const addLocalTask = (task) => {
    if (!taskMatchesFilters(task)) return;
    setTotalTasks((total) => total + 1);
    setTasks((prevTasks) => [task, ...prevTasks]);
  };

  const updateLocalTask = (updatedTask) => {
    setTasks((prevTasks) => {
      const filteredTasks = prevTasks.filter((task) => task.id !== updatedTask.id);
      if (!taskMatchesFilters(updatedTask)) return filteredTasks;
      return [updatedTask, ...filteredTasks];
    });
  };

  const deleteLocalTask = (id) => {
    setTasks((prevTasks) => prevTasks.filter((task) => task.id !== id));
    setTotalTasks((total) => Math.max(0, total - 1));
  };

  useEffect(() => {
    loadGroups().catch(handleError);
  }, []);

  useEffect(() => {
    loadTasks().catch(handleError);
  }, [activeGroup, status, page, query]);

  const stats = useMemo(
    () => ({
      total: tasks.length,
      done: tasks.filter((t) => t.status === 'completed').length,
    }),
    [tasks]
  );

  const totalPages = Math.max(1, Math.ceil(totalTasks / query.pageSize));
  const canGoPrevious = page > 1;
  const canGoNext = page < totalPages;

  const applyControls = () => {
    setPage(1);
    setQuery({ ...controls, pageSize: Number(controls.pageSize) });
  };

  function handleError(err) {
    setError(err.response?.data?.error || err.message || 'Something went wrong');
  }

  return (
    <div className="app-shell">
      <Navbar />
      <div className="dashboard">
        <GroupSidebar
          groups={groups}
          activeGroup={activeGroup}
          onSelect={(value) => {
            setActiveGroup(value);
            setPage(1);
          }}
          onCreate={(payload) =>
            groupApi.createGroup(payload).then(loadGroups).catch(handleError)
          }
          onShare={(id, payload) =>
            groupApi.shareGroup(id, payload).then(loadGroups).catch(handleError)
          }
          onDelete={(id) =>
            groupApi.deleteGroup(id).then(() => {
              if (activeGroup === id) {
                setActiveGroup(null);
              }
              loadGroups();
            }).catch(handleError)
          }
          onUpdate={(id, payload) =>
            groupApi.updateGroup(id, payload).then(loadGroups).catch(handleError)
          }
        />
        <main className="workspace">
          <section className="toolbar">
            <div>
              <h1>Tasks</h1>
              <span>
                {stats.done}/{stats.total} completed
              </span>
            </div>
            
            {/* HERE IS THE CHANGE PASSWORD BUTTON WRAPPED WITH THE DROPDOWN */}
            <div style={{ display: 'flex', gap: '15px', alignItems: 'center' }}>
              <select
                value={status}
                onChange={(e) => {
                  setStatus(e.target.value);
                  setPage(1);
                }}
              >
                <option value="">All</option>
                <option value="pending">Pending</option>
                <option value="completed">Completed</option>
              </select>
              
              <Link 
                to="/change-password" 
                style={{
                  padding: '6px 12px',
                  backgroundColor: '#f0f0f0',
                  color: '#333',
                  textDecoration: 'none',
                  borderRadius: '4px',
                  fontSize: '0.9rem',
                  border: '1px solid #ccc',
                  fontWeight: 'bold'
                }}
              >
                Change Password
              </Link>
            </div>
            
          </section>
          <section className="pagination-toolbar">
            <label>
              Page size
              <select
                value={controls.pageSize}
                onChange={(e) => setControls({ ...controls, pageSize: Number(e.target.value) })}
              >
                <option value={5}>5</option>
                <option value={10}>10</option>
                <option value={20}>20</option>
                <option value={50}>50</option>
              </select>
            </label>
            <label>
              Sort by
              <select
                value={controls.sortBy}
                onChange={(e) => setControls({ ...controls, sortBy: e.target.value })}
              >
                <option value="due_at">Deadline date</option>
                <option value="priority">Priority</option>
              </select>
            </label>
            <label>
              Order
              <select
                value={controls.sortOrder}
                onChange={(e) => setControls({ ...controls, sortOrder: e.target.value })}
              >
                <option value="asc">Ascending</option>
                <option value="desc">Descending</option>
              </select>
            </label>
            <button type="button" onClick={applyControls}>Load</button>
            <div className="page-controls">
              <button type="button" disabled={!canGoPrevious} onClick={() => setPage((p) => p - 1)}>
                Prev
              </button>
              <span>
                Page {page} of {totalPages} ({totalTasks} tasks)
              </span>
              <button type="button" disabled={!canGoNext} onClick={() => setPage((p) => p + 1)}>
                Next
              </button>
            </div>
          </section>
          {error && <div className="alert">{error}</div>}
          <TaskForm
            groups={groups}
            onCreate={(payload) =>
              taskApi
                .createTask(payload)
                .then(addLocalTask)
                .catch(handleError)
            }
          />
          <TaskList
            tasks={tasks}
            groups={groups}
            onToggle={(task) =>
              taskApi
                .updateTask(task.id, {
                  ...task,
                  status: task.status === 'completed' ? 'pending' : 'completed',
                })
                .then(updateLocalTask)
                .catch(handleError)
            }
            onUpdate={(updatedTask) =>
              taskApi
                .updateTask(updatedTask.id, updatedTask)
                .then(updateLocalTask)
                .catch(handleError)
            }
            onDelete={(id) =>
              taskApi.deleteTask(id).then(() => deleteLocalTask(id)).catch(handleError)
            }
            // onReminder={(payload) =>
            //   taskApi
            //     .createReminder(payload)
            //     .then(() => setError('Reminder saved'))
            //     .catch(handleError)
            // }
          />
        </main>
      </div>
    </div>
  );
}
