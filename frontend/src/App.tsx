import { useState } from 'react'

function App() {
  const [count, setCount] = useState(0)

  return (
    <div className="min-h-screen flex items-center justify-center bg-surface-base text-text-primary">
      <div className="text-center space-y-4">
        <h1 className="text-4xl font-bold font-mono text-brand-cta">MyQQBot</h1>
        <p className="text-text-secondary">Go + Wails v3 + OneBot + OpenAI</p>
        <button
          className="px-4 py-2 bg-brand-cta/20 text-brand-cta border border-brand-cta/30 rounded-lg hover:bg-brand-cta/30 transition-colors cursor-pointer"
          onClick={() => setCount(c => c + 1)}
        >
          Count: {count}
        </button>
      </div>
    </div>
  )
}

export default App
