package main

import (
	"github.com/1005281342/basic_component/ring_queue"
	"github.com/1005281342/test_tools/survey"
)

func main() {

	var (
		start = 100
		end   = 2000
		step  = 100
		cnt   = 10
	)
	// 并发写
	var rqb = ring_queue.NewRingQueueBlockRWLock(end)
	survey.RunIterations("ringQueueBlock_Insert", start, end, step,
		survey.Func2(func(x interface{}) { rqb.Insert(x) }, cnt))
	survey.RunIterations("ringQueueBlock_LPop", start, end, step,
		survey.Func2(func(x interface{}) { rqb.LPop() }, cnt))

	//var rq = ring_queue.NewRingQueueRWLock(end)
	//survey.RunIterations("ringQueue_Insert", start, end, step,
	//	survey.Func2(func(x interface{}) { rq.Insert(x) }, cnt))
	//survey.RunIterations("ringQueue_LPop", start, end, step,
	//	survey.Func2(func(x interface{}) { rq.LPop() }, cnt))

	// 准备数据
	var nums = survey.RandShuffle(end)
	for i := 0; i < end; i++ {
		//rq.Insert(nums[i])
		rqb.Insert(nums[i])
	}

	// 并发读
	survey.RunIterations("ringQueueBlock_Head", start, end, step,
		survey.Func2(func(x interface{}) { rqb.Head() }, cnt))
	survey.RunIterations("ringQueueBlock_Tail", start, end, step,
		survey.Func2(func(x interface{}) { rqb.Tail() }, cnt))
	survey.RunIterations("ringQueueBlock_Len", start, end, step,
		survey.Func2(func(x interface{}) { rqb.Len() }, cnt))
	survey.RunIterations("ringQueueBlock_IsFull", start, end, step,
		survey.Func2(func(x interface{}) { rqb.IsFull() }, cnt))
	survey.RunIterations("ringQueueBlock_Empty", start, end, step,
		survey.Func2(func(x interface{}) { rqb.Empty() }, cnt))

	//survey.RunIterations("ringQueue_Head", start, end, step,
	//	survey.Func2(func(x interface{}) { rq.Head() }, cnt))
	//survey.RunIterations("ringQueue_Tail", start, end, step,
	//	survey.Func2(func(x interface{}) { rq.Tail() }, cnt))
	//survey.RunIterations("ringQueue_Len", start, end, step,
	//	survey.Func2(func(x interface{}) { rq.Len() }, cnt))
	//survey.RunIterations("ringQueue_IsFull", start, end, step,
	//	survey.Func2(func(x interface{}) { rq.IsFull() }, cnt))
	//survey.RunIterations("ringQueue_Empty", start, end, step,
	//	survey.Func2(func(x interface{}) { rq.Empty() }, cnt))
}

/*
ringQueueBlock_Insert, 967, 1178, 1196, 802, 928, 877, 724, 746, 707, 645, 689, 665, 676, 641, 662, 680, 670, 632, 642, 666,
ringQueueBlock_LPop, 609, 762, 676, 716, 868, 725, 637, 672, 645, 618, 598, 640, 576, 682, 687, 635, 637, 723, 661, 660,
ringQueue_Insert, 691, 1178, 1268, 1793, 1259, 1357, 1174, 1173, 1276, 1380, 1344, 1159, 1017, 1074, 928, 945, 992, 927, 986, 949,
ringQueue_LPop, 543, 841, 1000, 1111, 1010, 887, 787, 764, 753, 725, 714, 766, 762, 768, 705, 685, 750, 753, 721, 656,
ringQueueBlock_Head, 581, 719, 657, 724, 707, 627, 726, 672, 738, 683, 749, 744, 687, 675, 668, 697, 681, 712, 1146, 861,
ringQueueBlock_Tail, 754, 754, 704, 782, 885, 832, 804, 1026, 978, 797, 826, 862, 995, 808, 826, 796, 836, 708, 743, 755,
ringQueueBlock_Len, 771, 678, 673, 783, 830, 884, 825, 899, 754, 819, 934, 863, 760, 848, 855, 842, 868, 870, 770, 774,
ringQueueBlock_IsFull, 787, 871, 799, 728, 834, 762, 728, 795, 843, 870, 775, 822, 853, 795, 766, 755, 760, 777, 1062, 999,
ringQueueBlock_Empty, 2816, 774, 838, 1026, 943, 891, 864, 786, 1035, 940, 939, 988, 549, 570, 713, 749, 859, 816, 776, 801,
ringQueue_Head, 739, 757, 856, 818, 747, 723, 824, 832, 858, 777, 996, 962, 950, 867, 865, 833, 806, 933, 884, 843,
ringQueue_Tail, 780, 965, 848, 830, 850, 700, 719, 777, 779, 921, 946, 1037, 745, 784, 738, 768, 730, 839, 804, 861,
ringQueue_Len, 1132, 919, 789, 704, 731, 815, 848, 967, 955, 947, 913, 874, 931, 909, 917, 840, 937, 789, 815, 740,
ringQueue_IsFull, 682, 580, 796, 775, 654, 699, 736, 697, 767, 822, 850, 880, 883, 925, 781, 903, 848, 827, 845, 760,
ringQueue_Empty, 622, 744, 711, 749, 763, 785, 950, 835, 751, 726, 826, 761, 730, 824, 813, 791, 833, 843, 868, 806,
*/
