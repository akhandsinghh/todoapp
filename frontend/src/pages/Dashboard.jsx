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

  const loadGroups = () => groupApi.listGroups().then(setGroups);
  const loadTasks = () =>
    taskApi.listTasks({ status, group_id: activeGroup || undefined }).then(setTasks);

  useEffect(() => {
    loadGroups().catch(handleError);
  }, []);

  useEffect(() => {
    loadTasks().catch(handleError);
  }, [activeGroup, status]);

  const stats = useMemo(
    () => ({
      total: tasks.length,
      done: tasks.filter((t) => t.status === 'completed').length,
    }),
    [tasks]
  );

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
          onSelect={setActiveGroup}
          onCreate={(payload) =>
            groupApi.createGroup(payload).then(loadGroups).catch(handleError)
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
              <select value={status} onChange={(e) => setStatus(e.target.value)}>
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
          {error && <div className="alert">{error}</div>}
          <TaskForm
            groups={groups}
            onCreate={(payload) =>
              taskApi.createTask(payload).then(loadTasks).catch(handleError)
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
                .then(loadTasks)
                .catch(handleError)
            }

            onUpdate={(updatedTask) =>
              taskApi
                .updateTask(updatedTask.id, updatedTask)
                .then(loadTasks)
                .catch(handleError)
            }
            
            onDelete={(id) => taskApi.deleteTask(id).then(loadTasks).catch(handleError)}
            onReminder={(payload) =>
              taskApi
                .createReminder(payload)
                .then(() => setError('Reminder saved'))
                .catch(handleError)
            }
          />
        </main>
      </div>
    </div>
  );
}