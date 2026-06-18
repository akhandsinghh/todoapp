import React, { useState } from 'react';

const colors = ['#4f46e5', '#059669', '#dc2626', '#d97706', '#0891b2'];

export default function GroupSidebar({ groups, activeGroup, onSelect, onCreate }) {
  const [name, setName] = useState('');
  const [color, setColor] = useState(colors[0]);

  const submit = (e) => {
    e.preventDefault();
    if (!name.trim()) return;
    onCreate({ name, color });
    setName('');
  };

  return (
    <aside className="sidebar">
      <button
        className={!activeGroup ? 'group active' : 'group'}
        onClick={() => onSelect(null)}
      >
        All tasks
      </button>
      {groups.map((g) => (
        <button
          key={g.id}
          className={activeGroup === g.id ? 'group active' : 'group'}
          onClick={() => onSelect(g.id)}
        >
          <i style={{ background: g.color }} />
          {g.name}
        </button>
      ))}
      <form className="group-form" onSubmit={submit}>
        <input
          placeholder="New group"
          value={name}
          onChange={(e) => setName(e.target.value)}
        />
        <div className="swatches">
          {colors.map((c) => (
            <button
              type="button"
              key={c}
              className={color === c ? 'selected' : ''}
              style={{ background: c }}
              onClick={() => setColor(c)}
              aria-label={c}
            />
          ))}
        </div>
        <button type="submit" 
        disabled={!name.trim()} 
        style={{
            opacity: !name.trim() ? 0.5 : 1,
            cursor: !name.trim() ? 'not-allowed' : 'pointer',
            transition: 'opacity 0.2s ease-in-out'
          }}>Add group</button>
      </form>
    </aside>
  );
}