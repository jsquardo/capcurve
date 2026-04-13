import type { ReactNode } from 'react'

interface Props {
  header: ReactNode  // full-width page title strip
  sidebar: ReactNode // filter panel
  main: ReactNode    // results area
}

// Two-column layout shell for the Playground page.
// header spans full width (owns its own background/border).
// On desktop: sidebar fixed-width left, main fills remaining space.
// On mobile: sidebar stacks above main with a bottom border instead of right border.
export default function PlaygroundLayout({ header, sidebar, main }: Props) {
  return (
    <div>
      {header}

      <div className="shell-container">
        <div className="flex min-h-[600px] flex-col lg:flex-row">
          <aside className="w-full shrink-0 border-b border-border py-6 lg:w-[280px] lg:border-b-0 lg:border-r lg:pr-6">
            {sidebar}
          </aside>
          <main className="min-w-0 flex-1 py-6 lg:pl-6">
            {main}
          </main>
        </div>
      </div>
    </div>
  )
}
