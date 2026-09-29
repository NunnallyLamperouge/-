import TablePanel from './TablePanel.vue'
import ImagePanel from './ImagePanel.vue'
import ListPanel from './ListPanel.vue'
import TextPanel from './TextPanel.vue'
import TitlePanel from './TitlePanel.vue'
import ReferencePanel from './ReferencePanel.vue'
import SchedulePanel from './SchedulePanel.vue'
import TocPanel from './TocPanel.vue'
import CoverTablePanel from './CoverTablePanel.vue'

export {
  TablePanel,
  ImagePanel,
  ListPanel,
  TextPanel,
  TitlePanel,
  ReferencePanel,
  SchedulePanel,
  TocPanel,
  CoverTablePanel
}

// 组件类型到面板组件的映射
export const componentPanelMap = {
  'table': TablePanel,
  'image': ImagePanel,
  'list': ListPanel,
  'text': TextPanel,
  'title': TitlePanel,
  'reference': ReferencePanel,
  'schedule': SchedulePanel,
  'toc': TocPanel,
  'cover-table': CoverTablePanel
} 