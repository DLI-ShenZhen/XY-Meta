### ***********************************************************************
### Author: Dehua Li <dehualiy@qq.com>;
### ***********************************************************************

#if ("BiocManager" %in% requirement_packages[,1]!=TRUE){
  #install.packages("getopt")
#}
requirement_packages=installed.packages()
if ("getopt" %in% requirement_packages[,1]!=TRUE){
  install.packages("getopt")
}
if ("xcms" %in% requirement_packages[,1]!=TRUE){
  BiocManager::install("xcms")
}
if ("MSnbase" %in% requirement_packages[,1]!=TRUE){
  BiocManager::install("MSnbase")
}

library(parallel)
library(getopt)
opt=data.frame(mass_tolerance=10,
               peak_width_min=20,
               peak_width_max=50,
               prefilter_min=3,
               prefilter_max=100,
               snthresh=4,
               mzdiff=0.001,
               noise=0,
               bw=5,
               minfrac=0.3,
               mzwid=0.015,
               group_max=1000,
               analysis_path="test",
               sample1="sample1",
               sample2="sample2",
               search="FALSE",
               xymeta_path="",
               paramater_path="")

spec = matrix(c(
  'mass_tolerance', 'p', 2, "double",
  'peak_width_min', 'i', 2, "double",
  'peak_width_max', 'I', 2, "double",
  'prefilter_min', 'a', 2, "double",
  'prefilter_max', 'A', 2, "double",
  'snthresh', 's', 2, "double",
  'mzdiff', 'm', 2, "double",
  'noise', 'o', 2, "double",
  'bw', 'b', 2, "double",
  'minfrac', 'r', 2, "double",
  'mzwid', 'z', 2, "double",
  'group_max', 'g', 2, "double",
  'analysis_path','W',2,"character",
  'sample1','T',2,"character",
  'sample2','N',2,"character",
  'search','S',2,"character",
  'xymeta_path','X',2,"character",
  'paramater_path','M',2,"character"), byrow=TRUE, ncol=4)

opts = getopt(spec)
if (is.null(opts$mass_tolerance)!=TRUE){
  opt$mass_tolerance=opts$mass_tolerance
}
if (is.null(opts$peak_width_min)!=TRUE){
  opt$peak_width_min=opts$peak_width_min
}
if (is.null(opts$peak_width_max)!=TRUE){
  opt$peak_width_max=opts$peak_width_max
}
if (is.null(opts$prefilter_min)!=TRUE){
  opt$prefilter_min=opts$prefilter_min
}
if (is.null(opts$prefilter_max)!=TRUE){
  opt$prefilter_max=opts$prefilter_max
}
if (is.null(opts$snthresh)!=TRUE){
  opt$snthresh=opts$snthresh
}
if (is.null(opts$mzdiff)!=TRUE){
  opt$mzdiff=opts$mzdiff
}
if (is.null(opts$noise)!=TRUE){
  opt$noise=opts$noise
}
if (is.null(opts$bw)!=TRUE){
  opt$bw=opts$bw
}
if (is.null(opts$minfrac)!=TRUE){
  opt$minfrac=opts$minfrac
}
if (is.null(opts$mzwid)!=TRUE){
  opt$mzwid=opts$mzwid
}
if (is.null(opts$group_max)!=TRUE){
  opt$group_max=opts$group_max
}
if (is.null(opts$analysis_path)!=TRUE){
  opts$analysis_path=gsub("\\\\","/",opts$analysis_path)
  opt$analysis_path=opts$analysis_path
}
if (is.null(opts$sample1)!=TRUE){
  opt$sample1=opts$sample1
}
if (is.null(opts$sample2)!=TRUE){
  opt$sample2=opts$sample2
}
if (is.null(opts$search)!=TRUE){
  opt$search=opts$search
}
if (is.null(opts$xymeta_path)!=TRUE){
  opts$xymeta_path=gsub("\\\\","/",opts$xymeta_path)
  opt$xymeta_path=opts$xymeta_path
}
if (is.null(opts$paramater_path)!=TRUE){
  opts$paramater_path=gsub("\\\\","/",opts$paramater_path)
  opt$paramater_path=opts$paramater_path
}

print(opt)
myAlign <- function (opt=opt) {
  
  ### +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
  ### These are the 4 variables you have to set for your own datafiles anything else runs automatically
  ### Set your working directory under Windows, where your netCDF files are stored
  ### Organize samples in subdirectories according to their class names WT/GM, Sick/Healthy etc.
  ### Important: use "/" not "\"
  myDir = opt$analysis_path
  myClass1 = opt$sample1
  myClass2 = opt$sample2
  myResultDir = "myAlign"
  ### +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
  
  ### change working directory to your files, see +++ section
  setwd(myDir)
  
  ### get working directory
  (WD <- getwd())
  
  ### load the xcms package
  library(xcms)
  
  ### you can get help typing the following command at the commandline
  # ?retcor
  
  ### finds peaks in NetCDF
  xset <- xcmsSet(method="centWave",ppm =opt$mass_tolerance,
                  peakwidth=c(opt$peak_width_min,
                              opt$peak_width_max),
                  snthresh=opt$snthresh,
                  prefilter=c(opt$prefilter_min, opt$prefilter_max),
                  mzCenterFun="wMean",integrate = 1L,
                  mzdiff = opt$mzdiff,
                  fitgauss = FALSE, noise = opt$noise,
                  firstBaselineCheck = TRUE, roiScales = numeric())
  
  ### print used files and memory usuage
  
  ### Group peaks together across samples and show fancy graphics
  ### you can remove the sleep timer (silent run) or set it to 0.001
  xset <- group(xset, sleep=0.0001)
  
  ### calculate retention time deviations for every time
  xset2 <- retcor(xset, family="s", plottype="m")
  
  ### Group peaks together across samples, set bandwitdh, change important m/z parameters here
  ### Syntax: group(object, bw = 30, minfrac = 0.5, minsamp= 1, mzwid = 0.25, max = 5, sleep = 0)
  xset2 <- group(xset2, bw =opt$bw,minfrac = opt$minfrac, mzwid = opt$mzwid, max = opt$group_max)
  
  ### identify peak groups and integrate samples
  xset3 <- fillPeaks(xset2)
  
  ### print statistics
  xset3
  
  ### output intensity matrix
  dat=groupval(xset3, "medret", "into")
  dat=rbind(group = as.character(phenoData(xset3)$class),dat)
  write.csv(dat, file="MyPeakTable.csv")
  
  ### create report and save the result in EXCEL file, print 20 important peaks as PNG
  ### reporttab <- diffreport(xset3, myClass1, myClass2, myResultDir, 20, metlin = 0.5)
  reporttab <- diffreport(xset3, myClass1, myClass2, myResultDir, metlin = 0.5)
  
  ### print file names
  # dir(path = ".", pattern = NULL, all.files = FALSE, full.names = FALSE, recursive = FALSE)
  
  ### output were done!
  print("Feature detection has been finished, open by yourself the file myAlign.tsv and pictures in myAlign_eic")
}

### gives CPU, system, TOTAL time in seconds
system.time(myAlign(opt))

### Currently R has no dual or multicore functionality, only a parallel library(snow)
### Benchmark on Dual Opteron 254 2.8 GHz with ARECA-1120 RAID5:
### 55 seconds total for the original 12 samples of the faahKO testset
### [1] 51.61 2.27 54.85 NA NA
### ***********************************************************************
### function finished
### ***********************************************************************

myDir = opt$analysis_path#work directory
myClass1 =paste(myDir,opt$sample1,sep = "/")#path of class1
myClass2 =paste(myDir,opt$sample2,sep = "/")#path of class2
quantify_data_path=Sys.glob(paste(myDir,"*myAlign.tsv",sep='/'))

Extra_sepctra=function(x){
  setwd(x)
  path_mzxml=Sys.glob(paste(x,"*.mzXML",sep='/'))#Obtain the path of mzXML
  cl=makeCluster(getOption("cl.cores",detectCores()))
  
  Extra_mz=function(x){
    
    library(MSnbase)
    name_lab=strsplit(x,"/")
    name_lab=strsplit(name_lab[[1]][length(name_lab[[1]])],".mz")
    msexp_ms2 <- readMSData(x,verbose = FALSE, centroided = TRUE)
    msexp_ms1 <- readMSData(x,verbose = FALSE, centroided = TRUE,msLevel. = 1)
    writeMgfData(msexp_ms2,paste(name_lab[[1]][1],"_ms2.mgf",sep = ""))
    writeMgfData(msexp_ms1,paste(name_lab[[1]][1],"_ms1.mgf",sep = ""))
    
    spectral_info=data.frame()
    for (i in 1:length(msexp_ms2)){
      test=data.frame(spectrum_scan=msexp_ms2[[i]]@precScanNum,
                      spectrum_index=msexp_ms2[[i]]@scanIndex,
                      precursormz=msexp_ms2[[i]]@precursorMz,
                      precursorintensity=msexp_ms2[[i]]@precursorIntensity,
                      precursorcharge=msexp_ms2[[i]]@precursorCharge,
                      energy=msexp_ms2[[i]]@collisionEnergy,
                      mslevel=msexp_ms2[[i]]@msLevel,
                      peakscount=msexp_ms2[[i]]@peaksCount,
                      precursorrt=msexp_ms2[[i]]@rt,
                      tic=msexp_ms2[[i]]@tic,
                      centroided=msexp_ms2[[i]]@centroided)
      spectral_info=rbind(spectral_info,test)
    }
    write.csv(spectral_info,paste(name_lab[[1]][1],"spectral_ms2_info.csv",sep = ""))
    
    spectral_info=data.frame()
    for (i in 1:length(msexp_ms1)){
      test=data.frame(spectrum_index=msexp_ms1[[i]]@scanIndex,
                      #precursormz=msexp_ms1[[i]]@precursorMz,
                      mslevel=msexp_ms1[[i]]@msLevel,
                      peakscount=msexp_ms1[[i]]@peaksCount,
                      precursorrt=msexp_ms1[[i]]@rt,
                      tic=msexp_ms1[[i]]@tic,
                      centroided=msexp_ms1[[i]]@centroided)
      spectral_info=rbind(spectral_info,test)
    }
    write.csv(spectral_info,paste(name_lab[[1]][1],"spectral_ms1_info.csv",sep = ""))
  }
  ###The kernel extra_sepatra function
  
  system.time({
    parLapply(cl,path_mzxml,Extra_mz)
  })
  ###Running program with mutiple
  stopCluster(cl)
}

Extra_sepctra(myClass1)
print("The mzXML files of sample1 have been extracted!")
Extra_sepctra(myClass2)
print("The mzXML files of sample1 have been extracted!")

###Database search program
SpecMatch=function(path_work,label,quantify_data_path,opt){
  setwd(path_work)
  #......Reading mgf filed......#
  path_mgf=Sys.glob(paste(getwd(),"*.mgf",sep='/'))
  command_line_input=paste("copy","*_ms2.mgf",paste(label,"_merge.mgf",sep = ""),sep = " ")
  shell(command_line_input)
  #......Running XY-Meta......#
  print("XY-Meta:Start to match")
  
  XY_Meta_path=opt$xymeta_path#path of search engine
  
  query_path=Sys.glob(paste(getwd(),"*_merge.mgf",sep='/'))#path of query mgf
  paramater_path=opt$paramater_path#path of paramater
  command_line_input=paste(XY_Meta_path,"-S",paramater_path[1],"-D",query_path,sep = " ")
  shell(command_line_input,wait=TRUE,intern = TRUE)
  
  #......Analyze the result from XY-Meta......
  #db.MS2_path=paste(getwd(),"/db.MS2.rData",sep = "")
  #load(db.MS2_path)
  print("Analyze match result")

  result_path=Sys.glob(paste(getwd(),"*.meta",sep='/'))
  #import the identification result of XY-Meta
  xyMetaResult=read.delim(result_path,header = T,sep = "\t")
  xyMetaResult=xyMetaResult[which(xyMetaResult$Score!=0),]
  flag=which(xyMetaResult$FDR<0.5)
  xyMetaResult_filter=xyMetaResult[1:flag[1],]#FDR filter
  rownames(xyMetaResult_filter)=xyMetaResult_filter[,1]
  
  #reading quantify data table
  quantify_data=read.delim(quantify_data_path,header = TRUE,sep = "\t",row.names = 1)
  annotation_label=c()
  #Annotating the match result
  for (i in 1:nrow(xyMetaResult_filter)){
    annotation_temp=""
    for (j in 1:nrow(quantify_data)){
      if (xyMetaResult$Query_precursor_mass[i]<=quantify_data$mzmax[j]&
          xyMetaResult$Query_precursor_mass[i]>=quantify_data$mzmin[j]&
          xyMetaResult$Query_precursor_retention_time[i]<=quantify_data$rtmax[j]&
          xyMetaResult$Query_precursor_retention_time[i]>=quantify_data$rtmin[j]){
        annotation_temp=paste(quantify_data$name[j],annotation_temp,sep = ";")
      }
    }
    annotation_label=append(annotation_label,annotation_temp)
  }
  xyMetaResult_filter=cbind(xyMetaResult_filter,MS1=annotation_label)
  
  #output the analysis result
  write.csv(xyMetaResult_filter,"Feature_identification_anno.csv",quote = TRUE)
  print("Match step has done")
}

if (opt$search=="TRUE"){
  SpecMatch(myClass1,opt$sample1,quantify_data_path,opt)#Database search
  SpecMatch(myClass2,opt$sample2,quantify_data_path,opt)#Database search
}

print("Finished")