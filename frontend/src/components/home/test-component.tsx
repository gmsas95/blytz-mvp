'use client'

import { useEffect, useState } from 'react'

export function TestComponent() {
  const [data, setData] = useState('Loading...')
  const [count, setCount] = useState(0)
  const [mounted, setMounted] = useState(false)

  useEffect(() => {
    console.log('TestComponent: useEffect triggered - MOUNTED')
    setMounted(true)

    // Test basic functionality
    const timer = setTimeout(() => {
      console.log('TestComponent: Timer fired, updating state')
      setData('React is working!')
      setCount(1)
    }, 100)

    return () => {
      console.log('TestComponent: Cleaning up timer')
      clearTimeout(timer)
    }
  }, [])

  useEffect(() => {
    console.log('TestComponent: State updated - data:', data, 'count:', count, 'mounted:', mounted)
  }, [data, count, mounted])

  // Force client-side rendering
  if (!mounted) {
    return (
      <div className="p-8 bg-yellow-100 border border-yellow-300 rounded-lg text-center">
        <h3 className="text-yellow-800 font-bold text-xl mb-2">Test Component (SSR)</h3>
        <p className="text-yellow-700">Mounting...</p>
      </div>
    )
  }

  return (
    <div className="p-8 bg-green-100 border border-green-300 rounded-lg text-center">
      <h3 className="text-green-800 font-bold text-xl mb-2">Test Component {count}</h3>
      <p className="text-green-700">{data}</p>
    </div>
  )
}