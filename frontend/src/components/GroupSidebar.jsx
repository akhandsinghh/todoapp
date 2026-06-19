import React, { useState } from 'react';

const colors = ['#4f46e5', '#059669', '#dc2626', '#d97706', '#0891b2'];

export default function GroupSidebar({ groups, activeGroup, onSelect, onCreate, onShare, onDelete }) {
  const [name, setName] = useState('');
  const [color, setColor] = useState(colors[0]);
  const [shareEmailByGroup, setShareEmailByGroup] = useState({});

  const submit = (e) => {
    e.preventDefault();
    if (!name.trim()) return;
    onCreate({ name, color });
    setName('');
  };

  const share = (e, groupId) => {
    e.preventDefault();
    const email = (shareEmailByGroup[groupId] || '').trim();
    if (!email) return;
    onShare(groupId, { email });
    setShareEmailByGroup((prev) => ({ ...prev, [groupId]: '' }));
  };

  return (
    <aside className="sidebar">
      <button
        className={!activeGroup ? 'group active' : 'group'}
        onClick={() => onSelect(null)}
      >
        All tasks
      </button>
      <button
        className={activeGroup === 'ungrouped' ? 'group active' : 'group'}
        onClick={() => onSelect('ungrouped')}
      >
        Ungrouped
      </button>
      {groups.map((g) => (
        <div key={g.id} className="group-item">
          <button
            className={activeGroup === g.id ? 'group active' : 'group'}
            onClick={() => onSelect(g.id)}
          >
            <i style={{ background: g.color }} />
            <span>{g.name}</span>
            <em>{g.role === 'creator' ? 'creator' : 'shared'}</em>
          </button>
          {g.role === 'creator' && (
            <>
              <form className="share-form" onSubmit={(e) => share(e, g.id)}>
                <input
                  type="email"
                  placeholder="Share by email"
                  value={shareEmailByGroup[g.id] || ''}
                  onChange={(e) =>
                    setShareEmailByGroup((prev) => ({ ...prev, [g.id]: e.target.value }))
                  }
                />
                <button type="submit" disabled={!(shareEmailByGroup[g.id] || '').trim()}>
                  Share
                </button>
              </form>
              <button
                className="delete-group-btn"
                onClick={() => {
                  if (window.confirm(`Delete group "${g.name}"? This action cannot be undone.`)) {
                    onDelete(g.id);
                  }
                }}
              >
                Delete
              </button>
            </>
          )}
        </div>
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
