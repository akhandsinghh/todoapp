import React, { useState } from 'react';

export default function TaskForm({ groups, onCreate }) {
  const [form, setForm] = useState({
    title: '',
    description: '',
    priority: 'medium',
    group_id: '',
    due_at: '',
  });

  const submit = (e) => {
    e.preventDefault();
    if (!form.title.trim()) return;
    onCreate({
      ...form,
      group_id: form.group_id ? Number(form.group_id) : null,
      due_at: form.due_at ? new Date(form.due_at).toISOString() : '',
    });
    setForm({ title: '', description: '', priority: 'medium', group_id: '', due_at: '' });
  };

  return (
    <form className="task-form" onSubmit={submit}>
      <input
        className="title-input"
        placeholder="Add a task"
        value={form.title}
        onChange={(e) => setForm({ ...form, title: e.target.value })}
      />
      <input
        placeholder="Notes"
        value={form.description}
        onChange={(e) => setForm({ ...form, description: e.target.value })}
      />
      <select
        value={form.priority}
        onChange={(e) => setForm({ ...form, priority: e.target.value })}
      >
        <option value="low">Low</option>
        <option value="medium">Medium</option>
        <option value="high">High</option>
      </select>
      <select
        value={form.group_id}
        onChange={(e) => setForm({ ...form, group_id: e.target.value })}
      >
        <option value="">No group</option>
        {groups.map((g) => (
          <option key={g.id} value={g.id}>
            {g.name}
          </option>
        ))}
      </select>
      <input
        type="datetime-local"
        value={form.due_at}
        onChange={(e) => setForm({ ...form, due_at: e.target.value })}
      />
      <button type="submit">Add</button>
    </form>
  );
}