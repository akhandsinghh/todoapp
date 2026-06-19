import React, { useState } from 'react';

const formatLocalDatetime = (dateStr) => {
  if (!dateStr) return '';
  const d = new Date(dateStr);
  const year = d.getFullYear();
  const month = String(d.getMonth() + 1).padStart(2, '0');
  const day = String(d.getDate()).padStart(2, '0');
  const hours = String(d.getHours()).padStart(2, '0');
  const minutes = String(d.getMinutes()).padStart(2, '0');
  return `${year}-${month}-${day}T${hours}:${minutes}`;
};

export default function TaskCard({ task, group, groups = [], onToggle, onDelete, onReminder, onUpdate }) {
  const [remindAt, setRemindAt] = useState(formatLocalDatetime(task.due_at));
  
  // New state variables for editing
  const [isEditing, setIsEditing] = useState(false);
  const [editForm, setEditForm] = useState({
    title: task.title || '',
    description: task.description || '',
    priority: task.priority || 'low',
    group_id: task.group_id || '',
    due_at: formatLocalDatetime(task.due_at) || ''
  });

  const done = task.status === 'completed';

  const submitReminder = (e) => {
    e.preventDefault();
    if (!remindAt) return;
    onReminder({
      task_id: task.id,
      remind_at: new Date(remindAt).toISOString(),
      message: task.title,
    });
  };

  const handleEditSubmit = (e) => {
    e.preventDefault();
    // Reformat local datetime to ISO standard for the Go backend
    const updatedTask = {
      ...task,
      ...editForm,
      due_at: editForm.due_at ? new Date(editForm.due_at).toISOString() : null,
      group_id: editForm.group_id === '' ? null : Number(editForm.group_id)
    };
    onUpdate(updatedTask);
    setIsEditing(false); // Close edit mode
  };

  // If in edit mode, render the form instead of the standard card
  if (isEditing) {
    return (
      <article className="task-card edit-mode">
        <form onSubmit={handleEditSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '10px', width: '100%' }}>
          <input
            type="text"
            required
            value={editForm.title}
            onChange={(e) => setEditForm({ ...editForm, title: e.target.value })}
            placeholder="Task Name"
          />
          <textarea
            value={editForm.description}
            onChange={(e) => setEditForm({ ...editForm, description: e.target.value })}
            placeholder="Description"
          />
          <div style={{ display: 'flex', gap: '10px', flexWrap: 'wrap' }}>
            <select 
              value={editForm.priority} 
              onChange={(e) => setEditForm({ ...editForm, priority: e.target.value })}
            >
              <option value="low">Low Priority</option>
              <option value="medium">Medium Priority</option>
              <option value="high">High Priority</option>
            </select>
            
            <select 
              value={editForm.group_id || ''} 
              onChange={(e) => setEditForm({ ...editForm, group_id: e.target.value })}
            >
              <option value="">No Group</option>
              {groups.map(g => (
                <option key={g.id} value={g.id}>{g.name}</option>
              ))}
            </select>
            
            <input
              type="datetime-local"
              value={editForm.due_at}
              onChange={(e) => setEditForm({ ...editForm, due_at: e.target.value })}
            />
          </div>
          <div style={{ display: 'flex', gap: '10px', marginTop: '10px' }}>
            <button type="submit">Save Changes</button>
            <button type="button" className="danger" onClick={() => setIsEditing(false)}>Cancel</button>
          </div>
        </form>
      </article>
    );
  }

  // Standard visual mode
  return (
    <article className={done ? 'task-card done' : 'task-card'}>
      <div className="task-main">
        <input type="checkbox" checked={done} onChange={() => onToggle(task)} />
        <div>
          <h3>{task.title}</h3>
          {task.description && <p>{task.description}</p>}
          <div className="meta">
            <span className={`priority ${task.priority}`}>{task.priority}</span>
            {group && <span>{group.name}</span>}
            {task.due_at && <span>{new Date(task.due_at).toLocaleString()}</span>}
          </div>
        </div>
      </div>
      <form className="reminder-row" onSubmit={submitReminder}>
        {/* <input
          type="datetime-local"
          value={remindAt}
          onChange={(e) => setRemindAt(e.target.value)}
        />
        <button type="submit">Remind</button> */}
        <button type="button" style={{ background: '#d97706', color: 'white' }} onClick={() => setIsEditing(true)}>
          Edit
        </button>
        <button type="button" className="danger" onClick={() => onDelete(task.id)}>
          Delete
        </button>
      </form>
    </article>
  );
}